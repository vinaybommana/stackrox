package matcher

import (
	clusterDataStore "github.com/stackrox/rox/central/cluster/datastore"
	deploymentMappings "github.com/stackrox/rox/central/deployment/mappings"
	imageMappings "github.com/stackrox/rox/central/image/mappings"
	processDataStore "github.com/stackrox/rox/central/processindicator/datastore"
	roleDataStore "github.com/stackrox/rox/central/rbac/k8srole/datastore"
	bindingDataStore "github.com/stackrox/rox/central/rbac/k8srolebinding/datastore"
	"github.com/stackrox/rox/central/searchbasedpolicies/builders"
	serviceAccountDataStore "github.com/stackrox/rox/central/serviceaccount/datastore"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	once sync.Once

	registry          Registry
	deploymentBuilder Builder
	imageBuilder      Builder
)

func intialize() {
	registry = NewRegistry(
		processDataStore.Singleton(),
		builders.K8sRBACQueryBuilder{
			Clusters:        clusterDataStore.Singleton(),
			K8sRoles:        roleDataStore.Singleton(),
			K8sBindings:     bindingDataStore.Singleton(),
			ServiceAccounts: serviceAccountDataStore.Singleton(),
		},
	)
	deploymentBuilder = NewBuilder(registry, deploymentMappings.OptionsMap)
	imageBuilder = NewBuilder(registry, imageMappings.OptionsMap)
}

// RegistrySingleton returns the registry used by the singleton matcher builders.
func RegistrySingleton() Registry {
	once.Do(intialize)
	return registry
}

// DeploymentBuilderSingleton Builder when you want to build Matchers for deployment policies.
func DeploymentBuilderSingleton() Builder {
	once.Do(intialize)
	return deploymentBuilder
}

// ImageBuilderSingleton Builder when you want to build Matchers for image policies.
func ImageBuilderSingleton() Builder {
	once.Do(intialize)
	return imageBuilder
}
