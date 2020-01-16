package sac

import (
	clusterDackBox "github.com/stackrox/rox/central/cluster/dackbox"
	deploymentDackBox "github.com/stackrox/rox/central/deployment/dackbox"
	globaldb "github.com/stackrox/rox/central/globaldb/dackbox"
	imageDackBox "github.com/stackrox/rox/central/image/dackbox"
	componentDackBox "github.com/stackrox/rox/central/imagecomponent/dackbox"
	namespaceDackBox "github.com/stackrox/rox/central/namespace/dackbox"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search/filtered"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	imageComponentSAC = sac.ForResource(resources.ImageComponent)

	imageComponentClusterPath = [][]byte{
		componentDackBox.Bucket,
		imageDackBox.Bucket,
		deploymentDackBox.Bucket,
		namespaceDackBox.Bucket,
		clusterDackBox.Bucket,
	}

	imageComponentNamespacePath = [][]byte{
		componentDackBox.Bucket,
		imageDackBox.Bucket,
		deploymentDackBox.Bucket,
		namespaceDackBox.Bucket,
	}

	// ImageComponentSACFilter represents SAC filter for image components
	ImageComponentSACFilter filtered.Filter
)

func init() {
	var err error
	ImageComponentSACFilter, err = filtered.NewSACFilter(
		filtered.WithResourceHelper(imageComponentSAC),
		filtered.WithGraphProvider(globaldb.GetGlobalDackBox()),
		filtered.WithClusterPath(imageComponentClusterPath...),
		filtered.WithNamespacePath(imageComponentNamespacePath...),
	)
	utils.Must(err)
}
