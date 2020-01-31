package converter

import (
	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/generated/storage"
)

// ProtoImageComponentToEmbeddedImageScanComponent converts a *storage.ImageComponent proto object to *storage.EmbeddedImageScanComponent proto object
// `vulns` and `layer_index` does not get set.
func ProtoImageComponentToEmbeddedImageScanComponent(component *storage.ImageComponent) *storage.EmbeddedImageScanComponent {
	return &storage.EmbeddedImageScanComponent{
		Name:     component.GetName(),
		Version:  component.GetVersion(),
		License:  proto.Clone(component.GetLicense()).(*storage.License),
		Priority: component.GetPriority(),
		Source:   component.GetSource(),
	}
}
