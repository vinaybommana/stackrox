package centralclient

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/installation/store"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/telemetry/phonehome"
	"github.com/stackrox/rox/pkg/version"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	apiPathsAnnotation = "rhacs.redhat.com/telemetry-apipaths"
	tenantIDLabel      = "rhacs.redhat.com/tenant"
)

var (
	config = &phonehome.Config{
		ClientID: "11102e5e-ca16-4f2b-8d2e-e9e04e8dc531",
	}
	once sync.Once
	log  = logging.LoggerForModule()
)

func getInstanceConfig() (*phonehome.Config, map[string]any, error) {
	rc, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, err
	}
	clientset, err := kubernetes.NewForConfig(rc)
	if err != nil {
		return nil, nil, err
	}
	v, err := clientset.ServerVersion()
	if err != nil {
		return nil, nil, err
	}

	deployments := clientset.AppsV1().Deployments(env.Namespace.Setting())
	central, err := deployments.Get(context.Background(), "central", v1.GetOptions{})
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot get central deployment")
	}

	paths := central.GetAnnotations()[apiPathsAnnotation]
	trackedPaths = set.NewFrozenSet(strings.Split(paths, ",")...)

	orchestrator := storage.ClusterType_KUBERNETES_CLUSTER.String()
	if env.OpenshiftAPI.BooleanSetting() {
		orchestrator = storage.ClusterType_OPENSHIFT_CLUSTER.String()
	}

	ii, _, err := store.Singleton().Get(sac.WithAllAccess(context.Background()))
	if err != nil || ii == nil {
		return nil, nil, errors.Wrap(err, "cannot get installation information")
	}
	centralID := ii.Id

	tenantID := central.GetLabels()[tenantIDLabel]
	// Consider on-prem central a tenant of itself:
	if tenantID == "" {
		tenantID = centralID
	}

	return &phonehome.Config{
			ClientID:   centralID,
			ClientName: "Central",
			GroupID:    tenantID,
		},
		map[string]any{
			"Central version":    version.GetMainVersion(),
			"Chart version":      version.GetChartVersion(),
			"Orchestrator":       orchestrator,
			"Kubernetes version": v.GitVersion,
		}, nil
}

// InstanceConfig collects the central instance telemetry configuration from
// central Deployment annotations and orchestrator properties. The collected
// data is used for instance identification.
func InstanceConfig() *phonehome.Config {
	once.Do(func() {
		if !phonehome.Enabled() {
			log.Info("Phonehome telemetry collection disabled.")
			return
		}
		cfg, props, err := getInstanceConfig()
		if err != nil {
			log.Errorf("Failed to get telemetry configuration: %v. Using hardcoded values.", err)
			return
		}
		config = cfg
		log.Info("Central ID: ", config.ClientID)
		log.Info("Tenant ID: ", config.GroupID)
		log.Info("API path telemetry enabled for: ", trackedPaths)

		for event, funcs := range interceptors {
			for _, f := range funcs {
				config.AddInterceptorFunc(event, f)
			}
		}
		// Central adds itself to the tenant group, adding its properties to the
		// group properties:
		config.Telemeter().Group(config.GroupID, config.ClientID, props)
		// Add the local admin user as well, with no extra group properties:
		config.Telemeter().Group(config.GroupID, config.HashUserID("admin", ""), nil)

		config.Gatherer().AddGatherer(func(ctx context.Context) (map[string]any, error) {
			return props, nil
		})
	})
	return config
}
