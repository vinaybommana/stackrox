package converter

import "github.com/stackrox/rox/generated/storage"

// ProtoImageComponentToEmbeddedImageScanComponent converts a *storage.ImageComponent proto object to *storage.EmbeddedImageScanComponent proto object
// `vulns` and `layer_index` does not get set.
func ProtoImageComponentToEmbeddedImageScanComponent(component *storage.ImageComponent) *storage.EmbeddedImageScanComponent {
	return &storage.EmbeddedImageScanComponent{
		Name:     component.GetName(),
		Version:  component.GetVersion(),
		License:  convertToEmbeddedLicense(component.GetLicense()),
		Priority: component.GetPriority(),
		Source:   convertSource(component.GetSource()),
		Location: component.GetLocation(),
	}
}

func convertToEmbeddedLicense(input *storage.ImageComponent_License) *storage.License {
	return &storage.License{
		Name: input.GetName(),
		Type: input.GetType(),
		Url:  input.GetUrl(),
	}
}

func convertSource(source storage.ImageComponent_SourceType) storage.EmbeddedImageScanComponent_SourceType {
	switch source {
	case storage.ImageComponent_OS:
		return storage.EmbeddedImageScanComponent_OS
	case storage.ImageComponent_PYTHON:
		return storage.EmbeddedImageScanComponent_PYTHON
	case storage.ImageComponent_JAVA:
		return storage.EmbeddedImageScanComponent_JAVA
	case storage.ImageComponent_RUBY:
		return storage.EmbeddedImageScanComponent_RUBY
	case storage.ImageComponent_NODEJS:
		return storage.EmbeddedImageScanComponent_NODEJS
	default:
		return storage.EmbeddedImageScanComponent_OS
	}
}
