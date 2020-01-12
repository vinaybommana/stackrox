package matcher

import (
	processDataStore "github.com/stackrox/rox/central/processindicator/datastore"
	"github.com/stackrox/rox/central/searchbasedpolicies"
	"github.com/stackrox/rox/central/searchbasedpolicies/builders"
	"github.com/stackrox/rox/central/searchbasedpolicies/fields"
	"github.com/stackrox/rox/pkg/logging"
)

var log = logging.LoggerForModule()

// Registry is the registry of top-level query builders.
// Policy evaluation is effectively a conjunction of these.
type Registry []searchbasedpolicies.PolicyQueryBuilder

// NewRegistry returns a new registry of the builders with the given underlying datastore for fetching process indicators.
func NewRegistry(processIndicators processDataStore.DataStore,
	k8sRBACBuilder builders.K8sRBACQueryBuilder) Registry {
	reg := []searchbasedpolicies.PolicyQueryBuilder{
		fields.ImageNameQueryBuilder,
		fields.ImageAgeQueryBuilder,
		builders.NewDockerFileLineQueryBuilder(),
		builders.CVSSQueryBuilder{},
		builders.CVEQueryBuilder{},
		fields.ComponentQueryBuilder,
		fields.DisallowedAnnotationQueryBuilder,
		fields.ScanAgeQueryBuilder,
		builders.ScanExistsQueryBuilder{},
		builders.EnvQueryBuilder{},
		fields.VolumeQueryBuilder,
		fields.PortQueryBuilder,
		fields.RequiredLabelQueryBuilder,
		fields.RequiredAnnotationQueryBuilder,
		builders.PrivilegedQueryBuilder{},
		builders.NewAddCapQueryBuilder(),
		builders.NewDropCapQueryBuilder(),
		fields.ResourcePolicy,
		builders.ProcessQueryBuilder{
			ProcessGetter: processIndicators,
		},
		builders.ReadOnlyRootFSQueryBuilder{},
		builders.PortExposureQueryBuilder{},
		builders.ProcessWhitelistingBuilder{},
		builders.HostMountQueryBuilder{},
		k8sRBACBuilder,
	}
	return reg
}
