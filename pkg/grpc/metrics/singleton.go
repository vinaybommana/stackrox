package metrics

import "github.com/stackrox/rox/pkg/sync"

var (
	metricsInit sync.Once

	metrics GRPCMetrics
)

// Singleton returns a singleton of a GRPCMetrics stuct
func Singleton() GRPCMetrics {
	metricsInit.Do(func() {
		metrics = newGRPCMetrics()
	})
	return metrics
}
