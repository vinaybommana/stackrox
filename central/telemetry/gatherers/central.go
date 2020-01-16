package gatherers

import (
	"github.com/stackrox/rox/pkg/telemetry/data"
	"github.com/stackrox/rox/pkg/telemetry/gatherers"
)

// CentralGatherer objects will gather and return telemetry information about this Central
type CentralGatherer struct {
	databaseGatherer      *databaseGatherer
	apiGatherer           *apiGatherer
	componentInfoGatherer *gatherers.ComponentInfoGatherer
}

// NewCentralGatherer creates and returns a CentralGatherer object
func NewCentralGatherer(databaseGatherer *databaseGatherer, apiGatherer *apiGatherer, componentInfoGatherer *gatherers.ComponentInfoGatherer) *CentralGatherer {
	return &CentralGatherer{
		databaseGatherer:      databaseGatherer,
		apiGatherer:           apiGatherer,
		componentInfoGatherer: componentInfoGatherer,
	}
}

// Gather returns telemetry information about this Central
func (c *CentralGatherer) Gather() *data.CentralInfo {
	centralComponent := &data.CentralInfo{
		RoxComponentInfo: c.componentInfoGatherer.Gather(),
		License:          nil,
		Storage:          c.databaseGatherer.Gather(),
		APIStats:         c.apiGatherer.Gather(),
	}
	return centralComponent
}
