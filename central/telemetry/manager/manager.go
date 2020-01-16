package manager

import (
	"context"

	"github.com/etcd-io/bbolt"
	"github.com/pkg/errors"
	licenseManager "github.com/stackrox/rox/central/license/manager"
	"github.com/stackrox/rox/central/telemetry/gatherers"
	"github.com/stackrox/rox/central/telemetry/manager/internal/store"
	"github.com/stackrox/rox/generated/storage"
)

// Manager manages telemetry configuration, collection, and sending.
type Manager interface {
	GetTelemetryConfig(ctx context.Context) (*storage.TelemetryConfiguration, error)
	UpdateTelemetryConfig(ctx context.Context, config *storage.TelemetryConfiguration) error
}

// NewManager creates a new telemetry manager. The manager starts running immediately, and keeps running until the
// given context expires.
func NewManager(ctx context.Context, boltDB *bbolt.DB, gatherer *gatherers.CentralGatherer, licenseMgr licenseManager.LicenseManager) (Manager, error) {
	telemetryStore, err := store.New(boltDB)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create telemetry store")
	}

	return newManager(ctx, telemetryStore, gatherer, licenseMgr), nil
}
