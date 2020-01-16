package sac

import (
	clusterDackBox "github.com/stackrox/rox/central/cluster/dackbox"
	cveDackBox "github.com/stackrox/rox/central/cve/dackbox"
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
	cveSAC = sac.ForResource(resources.CVE)

	cveClusterPath = [][]byte{
		cveDackBox.Bucket,
		componentDackBox.Bucket,
		imageDackBox.Bucket,
		deploymentDackBox.Bucket,
		namespaceDackBox.Bucket,
		clusterDackBox.Bucket,
	}

	cveNamespacePath = [][]byte{
		cveDackBox.Bucket,
		componentDackBox.Bucket,
		imageDackBox.Bucket,
		deploymentDackBox.Bucket,
		namespaceDackBox.Bucket,
	}

	// CVESACFilter represents the SAC filter for CVEs
	CVESACFilter filtered.Filter
)

func init() {
	var err error
	CVESACFilter, err = filtered.NewSACFilter(
		filtered.WithResourceHelper(cveSAC),
		filtered.WithGraphProvider(globaldb.GetGlobalDackBox()),
		filtered.WithClusterPath(cveClusterPath...),
		filtered.WithNamespacePath(cveNamespacePath...),
	)
	utils.Must(err)
}
