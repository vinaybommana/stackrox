package datastore

import (
	componentCVEEdgeIndexer "github.com/stackrox/rox/central/componentcveedge/index"
	cveIndexer "github.com/stackrox/rox/central/cve/index"
	globaldb "github.com/stackrox/rox/central/globaldb/dackbox"
	"github.com/stackrox/rox/central/globalindex"
	imageIndexer "github.com/stackrox/rox/central/image/index"
	componentIndexer "github.com/stackrox/rox/central/imagecomponent/index"
	"github.com/stackrox/rox/central/imagecomponent/search"
	"github.com/stackrox/rox/central/imagecomponent/store/dackbox"
	imageComponentEdgeIndexer "github.com/stackrox/rox/central/imagecomponentedge/index"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	once sync.Once

	ad DataStore
)

func initialize() {
	if !features.Dackbox.Enabled() {
		ad = nil
		return
	}
	storage, err := dackbox.New(globaldb.GetGlobalDackBox())
	utils.Must(err)

	indexer := componentIndexer.New(globalindex.GetGlobalIndex())
	searcher := search.New(storage, globaldb.GetGlobalDackBox(),
		cveIndexer.Singleton(),
		componentCVEEdgeIndexer.Singleton(),
		componentIndexer.Singleton(),
		imageComponentEdgeIndexer.Singleton(),
		imageIndexer.New(globalindex.GetGlobalIndex()))

	ad, err = New(storage, indexer, searcher)
	utils.Must(err)
}

// Singleton provides the interface for non-service external interaction.
func Singleton() DataStore {
	once.Do(initialize)
	return ad
}
