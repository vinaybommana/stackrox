package telemetry

import (
	"context"

	"github.com/stackrox/rox/central/sensor/service/common"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/pkg/concurrency"
)

// KubernetesInfoChunkCallback is a callback function that handles a single chunk of Kubernetes info returned from the sensor.
type KubernetesInfoChunkCallback func(ctx concurrency.ErrorWaitable, chunk *central.TelemetryResponsePayload_KubernetesInfo) error

// Controller handles requesting telemetry data from remote clusters.
type Controller interface {
	PullKubernetesInfo(ctx context.Context, cb KubernetesInfoChunkCallback) error
	ProcessTelemetryDataResponse(resp *central.PullTelemetryDataResponse) error
}

// NewController creates and returns a new controller for telemetry data.
func NewController(injector common.MessageInjector, stopSig concurrency.ReadOnlyErrorSignal) Controller {
	return newController(injector, stopSig)
}
