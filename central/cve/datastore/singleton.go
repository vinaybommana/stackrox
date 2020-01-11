package datastore

import (
	"github.com/stackrox/rox/central/cve/index"
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
	storage, err := dackbox.New(globaldb.GetGlobalDackBox())
	utils.Must(err)

	indexer := index.New(globalindex.GetGlobalIndex())
	searcher := search.New(storage, indexer)

	ds, err = New(storage, indexer, searcher)
	utils.Must(err)
}

// Singleton returns a singleton instance of cve datastore
func Singleton() DataStore {
	if !features.Dackbox.Enabled() {
		return nil
	}
	once.Do(initialize)
	return ds
}
