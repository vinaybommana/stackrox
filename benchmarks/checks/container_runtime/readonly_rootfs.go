package containerruntime

import (
	"github.com/stackrox/rox/benchmarks/checks/utils"
	"github.com/stackrox/rox/generated/storage"
)

type readonlyRootfsBenchmark struct{}

func (c *readonlyRootfsBenchmark) Definition() utils.Definition {
	return utils.Definition{
		BenchmarkCheckDefinition: storage.BenchmarkCheckDefinition{
			Name:        "CIS Docker v1.1.0 - 5.12",
			Description: "Ensure the container's root filesystem is mounted as read only",
		}, Dependencies: []utils.Dependency{utils.InitContainers},
	}
}

func (c *readonlyRootfsBenchmark) Run() (result storage.BenchmarkCheckResult) {
	utils.Pass(&result)
	for _, container := range utils.ContainersRunning {
		if !container.HostConfig.ReadonlyRootfs {
			utils.Warn(&result)
			utils.AddNotef(&result, "Container '%v' (%v) does not have a readonly rootfs", container.ID, container.Name)
		}
	}
	return
}

// NewReadonlyRootfsBenchmark implements CIS-5.12
func NewReadonlyRootfsBenchmark() utils.Check {
	return &readonlyRootfsBenchmark{}
}
