package gatherers

import (
	"github.com/stackrox/rox/central/license/manager"
	licenseproto "github.com/stackrox/rox/generated/shared/license"
	"github.com/stackrox/rox/pkg/telemetry/data"
	"github.com/stackrox/rox/pkg/telemetry/gatherers"
)

// CentralGatherer objects will gather and return telemetry information about this Central
type CentralGatherer struct {
	licenseMgr manager.LicenseManager

	databaseGatherer      *databaseGatherer
	apiGatherer           *apiGatherer
	componentInfoGatherer *gatherers.ComponentInfoGatherer
}

// NewCentralGatherer creates and returns a CentralGatherer object
func NewCentralGatherer(licenseMgr manager.LicenseManager, databaseGatherer *databaseGatherer, apiGatherer *apiGatherer, componentInfoGatherer *gatherers.ComponentInfoGatherer) *CentralGatherer {
	return &CentralGatherer{
		licenseMgr:            licenseMgr,
		databaseGatherer:      databaseGatherer,
		apiGatherer:           apiGatherer,
		componentInfoGatherer: componentInfoGatherer,
	}
}

// Gather returns telemetry information about this Central
func (c *CentralGatherer) Gather() *data.CentralInfo {
	var activeLicense *licenseproto.License
	if c.licenseMgr != nil {
		activeLicense = c.licenseMgr.GetActiveLicense()
	}
	centralComponent := &data.CentralInfo{
		RoxComponentInfo: c.componentInfoGatherer.Gather(),
		License:          (*data.LicenseJSON)(activeLicense),
		Storage:          c.databaseGatherer.Gather(),
		APIStats:         c.apiGatherer.Gather(),
	}
	return centralComponent
}
