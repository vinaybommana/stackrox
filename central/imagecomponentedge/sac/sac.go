package sac

import (
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search/filtered"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	imageSAC = sac.ForResource(resources.Image)

	// ImageComponentEdgeSACFilter represents the SAC filter for images-component edges
	ImageComponentEdgeSACFilter filtered.Filter
)

func init() {
	var err error
	ImageComponentEdgeSACFilter, err = filtered.NewSACFilter(
		filtered.WithResourceHelper(imageSAC),
	)
	utils.Must(err)
}
