package violations

import (
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/booleanpolicy/augmentedobjs"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/set"
)

// ContextQueryFields is a map of lifecycle stage to violation context fields.
type ContextQueryFields map[storage.LifecycleStage]set.FrozenStringSet

// Context Fields to be added to queries
var (
	// Build stage context fields
	ImageContextFields = newContextFields(
		[]string{search.ImageName.String()},
		[]string{search.ImageName.String(), augmentedobjs.ContainerNameCustomTag})
	EnvVarContextFields = newContextFields(
		nil,
		[]string{augmentedobjs.ContainerNameCustomTag})
	VulnContextFields = newContextFields(
		[]string{search.ImageName.String(), search.CVE.String(), search.CVSS.String(), augmentedobjs.ComponentAndVersionCustomTag},
		[]string{search.ImageName.String(), augmentedobjs.ContainerNameCustomTag, search.CVE.String(), search.CVSS.String(), augmentedobjs.ComponentAndVersionCustomTag})

	// Deploy stage context fields
	VolumeContextFields = newContextFields(
		[]string{},
		[]string{augmentedobjs.ContainerNameCustomTag, search.VolumeName.String(), search.VolumeSource.String(), search.VolumeDestination.String(), search.VolumeReadonly.String(), search.VolumeType.String()})
	ContainerContextFields = newContextFields(
		[]string{},
		[]string{augmentedobjs.ContainerNameCustomTag})
	ResourceContextFields = newContextFields(
		[]string{},
		[]string{augmentedobjs.ContainerNameCustomTag})
	// TODO(rc) can we add container name to port context?
	PortContextFields = newContextFields(
		[]string{},
		[]string{augmentedobjs.ContainerNameCustomTag, search.Port.String(), search.PortProtocol.String()})
)

func newContextFields(buildStageContext []string, deployStageContext []string) ContextQueryFields {
	return ContextQueryFields{
		storage.LifecycleStage_BUILD:  set.NewFrozenStringSet(buildStageContext...),
		storage.LifecycleStage_DEPLOY: set.NewFrozenStringSet(deployStageContext...),
	}
}
