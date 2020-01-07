package datastore

import (
	"context"

	"github.com/stackrox/rox/central/telemetry/datastore/internal/store"
	"github.com/stackrox/rox/generated/storage"
)

// DataStore manages the telemetry configuration
type DataStore interface {
	GetConfig(ctx context.Context) (*storage.TelemetryConfiguration, error)
	SetConfig(ctx context.Context, config *storage.TelemetryConfiguration) (*storage.TelemetryConfiguration, error)
}

// New returns a new instance of a DataStore
func New(store store.Store) DataStore {
	return &dataStoreImpl{
		store: store,
	}
}
