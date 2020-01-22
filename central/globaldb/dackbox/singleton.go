package dakcbox

import (
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/pkg/dackbox"
	"github.com/stackrox/rox/pkg/dackbox/indexer"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	// GraphBucket specifies the prefix for the id map DackBox tracks and stores in the DB.
	GraphBucket = []byte("dackbox_graph")

	registry indexer.IndexRegistry
	duckBox  *dackbox.DackBox

	dackBoxInit sync.Once

	log = logging.LoggerForModule()
)

// GetGlobalDackBox returns the global dackbox.DackBox instance.
func GetGlobalDackBox() *dackbox.DackBox {
	initializeDackBox()
	return duckBox
}

// GetIndexerRegistry returns the registry of indices that DackBox will use to index items in the queue.
func GetIndexerRegistry() indexer.IndexRegistry {
	initializeDackBox()
	return registry
}

func initializeDackBox() {
	dackBoxInit.Do(func() {
		if !features.Dackbox.Enabled() {
			return
		}

		globaldb.RegisterBucket(GraphBucket, "Graph Keys")

		registry = indexer.NewIndexRegistry()

		var err error
		duckBox, err = dackbox.NewDackBox(globaldb.GetGlobalBadgerDB(), GraphBucket)
		if err != nil {
			log.Panicf("Could not load stored indices: %v", err)
		}
	})
}
