package data

// ProcessMemInfo contains telemetry data about the resources used by a process
type ProcessMemInfo struct {
	CurrentAllocBytes   uint64
	CurrentAllocObjects uint64

	TotalAllocBytes   uint64
	TotalAllocObjects uint64

	SysMemBytes uint64

	NumGCs     uint32  `json:"numGCs,omitempty"`
	GCFraction float64 `json:"gcFraction,omitempty"`
}

// ProcessInfo contains telemetry data about a process
type ProcessInfo struct {
	Name          string
	NumGoroutines int `json:"numGoroutines,omitempty"`
	NumCPUs       int `json:"numCPUs,omitempty"`

	Memory *ProcessMemInfo `json:",omitempty"`
}
