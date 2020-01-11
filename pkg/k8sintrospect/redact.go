package k8sintrospect

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// RedactSecret removes sensitive secret data from a secret object, but retains information about which keys are
// present.
func RedactSecret(secret *unstructured.Unstructured) {
	annotations := secret.GetAnnotations()
	delete(annotations, "kubectl.kubernetes.io/last-applied-configuration")
	secret.SetAnnotations(annotations)
	dataMap, found, err := unstructured.NestedMap(secret.UnstructuredContent(), "data")
	if found && err == nil {
		redactedStringData := make(map[string]string, len(dataMap))
		for key := range dataMap {
			redactedStringData[key] = "***REDACTED***"
		}
		_ = unstructured.SetNestedStringMap(secret.UnstructuredContent(), redactedStringData, "stringData")
	}
	unstructured.RemoveNestedField(secret.UnstructuredContent(), "data")
}
