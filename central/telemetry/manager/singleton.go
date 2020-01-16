package manager

import (
	"context"

	"github.com/stackrox/rox/central/globaldb"
	licenseSingletons "github.com/stackrox/rox/central/license/singleton"
	"github.com/stackrox/rox/central/telemetry/gatherers"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	instance     Manager
	instanceInit sync.Once
)

// Singleton returns the license manager singleton instance.
func Singleton() Manager {
	instanceInit.Do(func() {
		var err error
		instance, err = NewManager(context.Background(), globaldb.GetGlobalDB(), gatherers.Singleton(), licenseSingletons.ManagerSingleton())
		if err != nil {
			panic(err)
		}
	})
	return instance
}
