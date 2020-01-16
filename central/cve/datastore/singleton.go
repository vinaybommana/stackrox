package datastore

import (
	componentCVEEdgeIndexer "github.com/stackrox/rox/central/componentcveedge/index"
	cveIndexer "github.com/stackrox/rox/central/cve/index"
	imageIndexer "github.com/stackrox/rox/central/image/index"
	componentIndexer "github.com/stackrox/rox/central/imagecomponent/index"
	imageComponentEdgeIndexer "github.com/stackrox/rox/central/imagecomponentedge/index"

	"github.com/stackrox/rox/central/cve/search"
	"github.com/stackrox/rox/central/cve/store/dackbox"
	globaldb "github.com/stackrox/rox/central/globaldb/dackbox"
	"github.com/stackrox/rox/central/globalindex"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	once sync.Once

	ds DataStore
)

func initialize() {
	if !features.Dackbox.Enabled() {
		ds = nil
		return
	}
	storage, err := dackbox.New(globaldb.GetGlobalDackBox())
	utils.Must(err)

	searcher := search.New(storage, globaldb.GetGlobalDackBox(),
		cveIndexer.Singleton(),
		componentCVEEdgeIndexer.Singleton(),
		componentIndexer.Singleton(),
		imageComponentEdgeIndexer.Singleton(),
		imageIndexer.New(globalindex.GetGlobalIndex()))

	ds, err = New(storage, cveIndexer.Singleton(), searcher)
	utils.Must(err)
}

// Singleton returns a singleton instance of cve datastore
func Singleton() DataStore {
	once.Do(initialize)
	return ds
}
