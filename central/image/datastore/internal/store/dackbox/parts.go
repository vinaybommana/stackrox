package dackbox

import (
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
)

var (
	log = logging.LoggerForModule()
)

type imageParts struct {
	image     *storage.Image
	listImage *storage.ListImage

	children []componentParts
}

type componentParts struct {
	edge      *storage.ImageComponentEdge
	component *storage.ImageComponent

	children []cveParts
}

type cveParts struct {
	edge *storage.ComponentCVEEdge
	cve  *storage.CVE
}
