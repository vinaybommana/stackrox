package gatherers

import (
	"os"
	"runtime"

	"github.com/stackrox/rox/pkg/telemetry/data"
	"github.com/stackrox/rox/pkg/version"
)

// ComponentInfoGatherer gathers generic information about a StackRox component(Centra, Scanner, etc...)
type ComponentInfoGatherer struct {
}

// NewComponentInfoGatherer creates and returns a ComponentInfoGatherer
func NewComponentInfoGatherer() *ComponentInfoGatherer {
	return &ComponentInfoGatherer{}
}

// Gather returns generic telemetry information about a StackRox component (Central, Scanner, etc...)
func (c *ComponentInfoGatherer) Gather() *data.RoxComponentInfo {
	return &data.RoxComponentInfo{
		Version:  version.GetMainVersion(),
		Process:  getProcessInfo(),
		Restarts: 0, //TODO: Figure out how to ge number of restarts
	}
}

func getProcessInfo() *data.ProcessInfo {
	return &data.ProcessInfo{
		Name:          os.Args[0],
		NumGoroutines: runtime.NumGoroutine(),
		NumCPUs:       runtime.NumCPU(),
		Memory:        getMemInfo(),
	}
}

func getMemInfo() *data.ProcessMemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &data.ProcessMemInfo{
		CurrentAllocBytes:   m.Alloc,
		CurrentAllocObjects: m.HeapObjects,
		TotalAllocBytes:     m.TotalAlloc,
		TotalAllocObjects:   m.Mallocs,
		SysMemBytes:         m.Sys,
		NumGCs:              m.NumGC,
		GCFraction:          m.GCCPUFraction,
	}
}
