package gatherers

import (
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/license/singleton"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/telemetry/gatherers"
)

var (
	gatherer     *CentralGatherer
	gathererInit sync.Once
)

// Singleton initializes and returns a CentralGatherer singleton
func Singleton() *CentralGatherer {
	gathererInit.Do(func() {
		gatherer = NewCentralGatherer(singleton.ManagerSingleton(), newDatabaseGatherer(newBadgerGatherer(globaldb.GetGlobalBadgerDB()), newBoltGatherer(globaldb.GetGlobalDB())), newAPIGatherer(), gatherers.NewComponentInfoGatherer())
	})
	return gatherer
}
