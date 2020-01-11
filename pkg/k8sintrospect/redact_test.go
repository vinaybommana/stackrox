package k8sintrospect

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestRedactSecret(t *testing.T) {
	secret := v1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "stackrox",
			Name:      "super-secret-data",
			Annotations: map[string]string{
				"kubectl.kubernetes.io/last-applied-configuration": "some config string that contains raw secret data",
			},
		},
		Data: map[string][]byte{
			"key.pem":  []byte("secret key data"),
			"cert.pem": []byte("not so secret cert data"),
		},
	}

	var unstructuredSecret unstructured.Unstructured
	require.NoError(t, scheme.Scheme.Convert(&secret, &unstructuredSecret, nil))

	RedactSecret(&unstructuredSecret)

	var redactedSecret v1.Secret
	require.NoError(t, scheme.Scheme.Convert(&unstructuredSecret, &redactedSecret, nil))

	expectedStringData := map[string]string{
		"cert.pem": "***REDACTED***",
		"key.pem":  "***REDACTED***",
	}

	assert.Equal(t, redactedSecret.StringData, expectedStringData)
	assert.Empty(t, redactedSecret.Data)
	assert.NotContains(t, redactedSecret.GetAnnotations(), "kubectl.kubernetes.io/last-applied-configuration")
}
