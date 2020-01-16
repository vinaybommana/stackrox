package sac

import (
	clusterDackBox "github.com/stackrox/rox/central/cluster/dackbox"
	deploymentDackBox "github.com/stackrox/rox/central/deployment/dackbox"
	globaldb "github.com/stackrox/rox/central/globaldb/dackbox"
	imageDackBox "github.com/stackrox/rox/central/image/dackbox"
	namespaceDackBox "github.com/stackrox/rox/central/namespace/dackbox"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search/filtered"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	imageSAC = sac.ForResource(resources.Image)

	imageClusterPath = [][]byte{
		imageDackBox.Bucket,
		deploymentDackBox.Bucket,
		namespaceDackBox.Bucket,
		clusterDackBox.Bucket,
	}

	imageNamespacePath = [][]byte{
		imageDackBox.Bucket,
		deploymentDackBox.Bucket,
		namespaceDackBox.Bucket,
	}

	// ImageSACFilter represents the SAC filter for images
	ImageSACFilter filtered.Filter
)

func init() {
	var err error
	ImageSACFilter, err = filtered.NewSACFilter(
		filtered.WithResourceHelper(imageSAC),
		filtered.WithGraphProvider(globaldb.GetGlobalDackBox()),
		filtered.WithClusterPath(imageClusterPath...),
		filtered.WithNamespacePath(imageNamespacePath...),
	)
	utils.Must(err)
}
