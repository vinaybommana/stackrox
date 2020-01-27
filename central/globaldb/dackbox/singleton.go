package dakcbox

import (
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/dackbox"
	"github.com/stackrox/rox/pkg/dackbox/indexer"
	"github.com/stackrox/rox/pkg/dackbox/utils/queue"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	// GraphBucket specifies the prefix for the id map DackBox tracks and stores in the DB.
	GraphBucket = []byte("dackbox_graph")
	// DirtyBucket specifies the prefix for the set of dirty keys (need re-indexing) to add to dackbox.
	DirtyBucket = []byte("dackbox_dirty")
	// ValidBucket specifies the prefix for storing validation information in dackbox.
	ValidBucket = []byte("dackbox_valid")

	toIndex       queue.WaitableQueue
	registry      indexer.IndexRegistry
	lazy          indexer.Lazy
	duckBox       *dackbox.DackBox
	globalKeyLock concurrency.KeyFence

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

// GetIndexQueue returns the queue of items waiting to be indexed.
func GetIndexQueue() queue.WaitableQueue {
	initializeDackBox()
	return toIndex
}

// GetKeyFence returns the global key fence.
func GetKeyFence() concurrency.KeyFence {
	initializeDackBox()
	return globalKeyLock
}

func initializeDackBox() {
	dackBoxInit.Do(func() {
		if !features.Dackbox.Enabled() {
			return
		}

		globaldb.RegisterBucket(GraphBucket, "Graph Keys")
		globaldb.RegisterBucket(DirtyBucket, "Dirty Keys")
		globaldb.RegisterBucket(ValidBucket, "Valid DackBox State")

		toIndex = queue.NewWaitableQueue(queue.NewQueue())
		registry = indexer.NewIndexRegistry()
		globalKeyLock = concurrency.NewKeyFence()

		var err error
		duckBox, err = dackbox.NewDackBox(globaldb.GetGlobalBadgerDB(), toIndex, GraphBucket, DirtyBucket, ValidBucket)
		if err != nil {
			log.Panicf("Could not load stored indices: %v", err)
		}

		lazy = indexer.NewLazy(toIndex, registry, duckBox.AckIndexed)
		lazy.Start()
	})
}
