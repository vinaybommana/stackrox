package data

// CollectorInfo contains telemetry data specific to StackRox' Collector sidecar
type CollectorInfo struct {
	Version string
	*RoxComponentInfo
}
