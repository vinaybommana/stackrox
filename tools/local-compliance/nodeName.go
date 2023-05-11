package main

import (
	"os"

	"github.com/stackrox/rox/pkg/orchestrators"
	"github.com/stackrox/rox/pkg/sync"
)

type nodeNameProvider interface {
	getNode() string
}

type dummyNodeNameProvider struct{}

func (dnp *dummyNodeNameProvider) getNode() string {
	return "Foo"
}

type envNodeNameProvider struct {
	once sync.Once
}

func (np *envNodeNameProvider) getNode() string {
	var node string
	np.once.Do(func() {
		node = os.Getenv(string(orchestrators.NodeName))
		if node == "" {
			log.Fatal("No node name found in the environment")
		}
	})
	return node
}
