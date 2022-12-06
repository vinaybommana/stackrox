package phonehome

import (
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/telemetry/phonehome/segment"
)

var (
	telemeter     Telemeter
	onceTelemeter sync.Once
)

// Enabled returns true if telemetry data collection is enabled.
func Enabled() bool {
	return segment.Enabled()
}

// TelemeterSingleton returns the instance of the telemeter.
func TelemeterSingleton() Telemeter {
	onceTelemeter.Do(func() {
		cfg := InstanceConfig()
		telemeter = segment.NewTelemeter(cfg.CentralID, cfg.Properties)
		// Central adds itself to the tenant group, adding its properties to the
		// group properties:
		telemeter.Group(cfg.GetGroupID(), cfg.CentralID, cfg.Properties)
		// Add the local admin user as well:
		telemeter.Group(cfg.GetGroupID(), "local:"+cfg.CentralID+":admin", cfg.Properties)
	})
	return telemeter
}