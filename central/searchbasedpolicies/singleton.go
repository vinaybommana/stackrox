package searchbasedpolicies

import (
	clusterDataStore "github.com/stackrox/rox/central/cluster/datastore"
	processDataStore "github.com/stackrox/rox/central/processindicator/datastore"
	roleDataStore "github.com/stackrox/rox/central/rbac/k8srole/datastore"
	bindingDataStore "github.com/stackrox/rox/central/rbac/k8srolebinding/datastore"
	k8sBuilder "github.com/stackrox/rox/central/searchbasedpolicies/builders"
	serviceAccountDataStore "github.com/stackrox/rox/central/serviceaccount/datastore"
	"github.com/stackrox/rox/pkg/search/options/deployments"
	"github.com/stackrox/rox/pkg/search/options/images"
	"github.com/stackrox/rox/pkg/searchbasedpolicies/matcher"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	once sync.Once

	registry          matcher.Registry
	deploymentBuilder matcher.Builder
	imageBuilder      matcher.Builder
)

func initialize() {
	registry = matcher.NewRegistry(
		processDataStore.Singleton(),
		k8sBuilder.K8sRBACQueryBuilder{
			Clusters:        clusterDataStore.Singleton(),
			K8sRoles:        roleDataStore.Singleton(),
			K8sBindings:     bindingDataStore.Singleton(),
			ServiceAccounts: serviceAccountDataStore.Singleton(),
		},
	)
	deploymentBuilder = matcher.NewBuilder(registry, deployments.OptionsMap)
	imageBuilder = matcher.NewBuilder(registry, images.OptionsMap)
}

// RegistrySingleton returns the registry used by the singleton matcher builders.
func RegistrySingleton() matcher.Registry {
	once.Do(initialize)
	return registry
}

// DeploymentBuilderSingleton Builder when you want to build Matchers for deployment policies.
func DeploymentBuilderSingleton() matcher.Builder {
	once.Do(initialize)
	return deploymentBuilder
}

// ImageBuilderSingleton Builder when you want to build Matchers for image policies.
func ImageBuilderSingleton() matcher.Builder {
	once.Do(initialize)
	return imageBuilder
}
