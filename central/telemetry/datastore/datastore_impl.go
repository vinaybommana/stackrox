package datastore

import (
	"context"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/central/telemetry/datastore/internal/store"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/sac"
)

type dataStoreImpl struct {
	store store.Store
}

var (
	telemetrySAC = sac.ForResource(resources.DebugLogs)
)

func (d *dataStoreImpl) GetConfig(ctx context.Context) (*storage.TelemetryConfiguration, error) {
	if ok, err := telemetrySAC.ScopeChecker(ctx, storage.Access_READ_ACCESS).Allowed(ctx); err != nil {
		return nil, err
	} else if !ok {
		// Usually we return nil, nil for permission denied so we don't leak existence information but in this case
		// there is only one config object and it always exists
		return nil, errors.New("permission denied")
	}

	config, err := d.store.GetTelemetryConfig()
	if err != nil {
		return nil, err
	}
	if config == nil {
		config = &storage.TelemetryConfiguration{
			Enabled: env.InitialTelemetryEnabledEnv.BooleanSetting(),
		}
		err := d.store.SetTelemetryConfig(config)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

func (d *dataStoreImpl) SetConfig(ctx context.Context, config *storage.TelemetryConfiguration) (*storage.TelemetryConfiguration, error) {
	if ok, err := telemetrySAC.ScopeChecker(ctx, storage.Access_READ_WRITE_ACCESS).Allowed(ctx); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("permission denied")
	}

	if err := d.store.SetTelemetryConfig(config); err != nil {
		return nil, err
	}
	return config, nil
}
