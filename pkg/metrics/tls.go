package metrics

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/fileutils"
	"github.com/stackrox/rox/pkg/mtls/certwatch"
	"github.com/stackrox/rox/pkg/mtls/verifier"
	"github.com/stackrox/rox/pkg/sync"
	"go.uber.org/zap"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type tlsConfigLoader struct {
	certDir           string
	clientCAConfigMap string
	clientCANamespace string
	k8sClient         *kubernetes.Clientset
	tlsConfig         *tls.Config

	cfgMutex sync.RWMutex
}

// NewTLSConfigLoader creates a new TLS config loader.
func NewTLSConfigLoader(certDir, clientCANamespace, clientCAConfigMap string) (*tlsConfigLoader, error) {
	tlsRootConfig := verifier.DefaultTLSServerConfig(nil, nil)
	tlsRootConfig.ClientAuth = tls.RequireAndVerifyClientCert
	loader := &tlsConfigLoader{
		certDir:           certDir,
		clientCANamespace: clientCANamespace,
		clientCAConfigMap: clientCAConfigMap,
		tlsConfig:         tlsRootConfig,
	}
	loader.tlsConfig.GetConfigForClient = loader.getClientConfigFunc()

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	loader.k8sClient = clientset
	return loader, nil
}

// WatchForChanges watches for changes of the server TLS certificate files and the client CA config map.
func (t *tlsConfigLoader) WatchForChanges() {
	// Watch for changes of server TLS certificate.
	certwatch.WatchCertDir(t.certDir, t.getCertificateFromDirectory, t.updateCertificate)

	// Watch for changes of client CA.
	go t.watchForClientCAChanges()
}

func (t *tlsConfigLoader) TLSConfig() *tls.Config {
	if t == nil {
		return nil
	}
	return t.tlsConfig
}

func (t *tlsConfigLoader) getClientConfigFunc() func(*tls.ClientHelloInfo) (*tls.Config, error) {
	return func(clientHello *tls.ClientHelloInfo) (*tls.Config, error) {
		t.cfgMutex.RLock()
		defer t.cfgMutex.RUnlock()
		return t.tlsConfig, nil
	}
}

func (t *tlsConfigLoader) getCertificateFromDirectory(dir string) (*tls.Certificate, error) {
	certFile := filepath.Join(dir, env.TLSCertFileName)
	if exists, err := fileutils.Exists(certFile); err != nil || !exists {
		if err != nil {
			log.Warnw("Error checking if monitoring TLS certificate file exists", zap.Error(err))
			return nil, err
		}
		log.Infof("Monitoring TLS certificate file %q does not exist. Skipping", certFile)
		return nil, nil
	}

	keyFile := filepath.Join(dir, env.TLSKeyFileName)
	if exists, err := fileutils.Exists(keyFile); err != nil || !exists {
		if err != nil {
			log.Warnw("Error checking if monitoring TLS key file exists", zap.Error(err))
			return nil, err
		}
		log.Infof("Monitoring TLS key file %q does not exist. Skipping", keyFile)
		return nil, nil
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, errors.Wrap(err, "loading monitoring certificate failed")
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, errors.Wrap(err, "parsing leaf certificate failed")
	}
	return &cert, nil
}

func (t *tlsConfigLoader) updateCertificate(cert *tls.Certificate) {
	t.cfgMutex.Lock()
	defer t.cfgMutex.Unlock()
	t.tlsConfig.Certificates = []tls.Certificate{*cert}
}

func (t *tlsConfigLoader) watchForClientCAChanges() {
	for {
		watcher, err := t.k8sClient.CoreV1().ConfigMaps(t.clientCANamespace).Watch(
			context.Background(),
			metav1.SingleObject(metav1.ObjectMeta{
				Name: t.clientCAConfigMap, Namespace: t.clientCANamespace,
			}))
		if err != nil {
			log.Errorw("Unable to create client CA watcher", zap.Error(err))
			continue
		}
		t.updateClientCA(watcher.ResultChan())
	}
}

func (t *tlsConfigLoader) updateClientCA(eventChannel <-chan watch.Event) {
	for {
		event, open := <-eventChannel
		if open {
			switch event.Type {
			case watch.Added:
				fallthrough
			case watch.Modified:
				if cm, ok := event.Object.(*v1.ConfigMap); ok {
					if caFile, ok := cm.Data["client-ca-file"]; ok {
						caCert := []byte(caFile)
						caCertPool := x509.NewCertPool()
						caCertPool.AppendCertsFromPEM(caCert)

						t.cfgMutex.Lock()
						t.tlsConfig.ClientCAs = caCertPool
						t.cfgMutex.Unlock()
					}
				}
			default:
			}
		} else {
			// If eventChannel is closed the server has closed the connection.
			// We want to return and create another watcher.
			return
		}
	}
}
