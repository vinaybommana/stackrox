package sac

import (
	clusterDackBox "github.com/stackrox/rox/central/cluster/dackbox"
	deploymentDackBox "github.com/stackrox/rox/central/deployment/dackbox"
	globaldb "github.com/stackrox/rox/central/globaldb/dackbox"
	namespaceDackBox "github.com/stackrox/rox/central/namespace/dackbox"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search/filtered"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	deploymentSAC = sac.ForResource(resources.Image)

	deploymentClusterPath = [][]byte{
		deploymentDackBox.Bucket,
		namespaceDackBox.Bucket,
		clusterDackBox.Bucket,
	}

	deploymentNamespacePath = [][]byte{
		deploymentDackBox.Bucket,
		namespaceDackBox.Bucket,
	}

	imageSACFilter filtered.Filter
	once           sync.Once
)

// GetSACFilter returns the sac filter for image ids.
func GetSACFilter() filtered.Filter {
	once.Do(func() {
		var err error
		imageSACFilter, err = filtered.NewSACFilter(
			filtered.WithResourceHelper(deploymentSAC),
			filtered.WithGraphProvider(globaldb.GetGlobalDackBox()),
			filtered.WithClusterPath(deploymentClusterPath...),
			filtered.WithNamespacePath(deploymentNamespacePath...),
		)
		utils.Must(err)
	})
	return imageSACFilter
}
