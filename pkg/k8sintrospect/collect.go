package k8sintrospect

import (
	"context"

	"k8s.io/client-go/rest"
)

// File is a file emitted by the K8s introspection feature.
type File struct {
	Path     string
	Contents []byte
}

// Collect collects Kubernetes data relevant to our deployment.
func Collect(ctx context.Context, collectionCfg Config, k8sClientConfig *rest.Config, filesC chan<- File) error {
	c, err := newCollector(ctx, k8sClientConfig, collectionCfg, filesC)
	if err != nil {
		return err
	}
	return c.Run()
}
