package datastore

import (
	"context"

	"github.com/blevesearch/bleve"
	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
	componentCVEEdgeIndexer "github.com/stackrox/rox/central/componentcveedge/index"
	cveIndexer "github.com/stackrox/rox/central/cve/index"
	globalDackBox "github.com/stackrox/rox/central/globaldb/dackbox"
	"github.com/stackrox/rox/central/image/datastore/internal/search"
	"github.com/stackrox/rox/central/image/datastore/internal/store"
	badgerStore "github.com/stackrox/rox/central/image/datastore/internal/store/badger"
	imageIndexer "github.com/stackrox/rox/central/image/index"
	componentIndexer "github.com/stackrox/rox/central/imagecomponent/index"
	imageComponentEdgeIndexer "github.com/stackrox/rox/central/imagecomponentedge/index"
	riskDS "github.com/stackrox/rox/central/risk/datastore"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	searchPkg "github.com/stackrox/rox/pkg/search"
)

// DataStore is an intermediary to AlertStorage.
//go:generate mockgen-wrapper DataStore
type DataStore interface {
	SearchListImages(ctx context.Context, q *v1.Query) ([]*storage.ListImage, error)
	ListImage(ctx context.Context, sha string) (*storage.ListImage, bool, error)

	Search(ctx context.Context, q *v1.Query) ([]searchPkg.Result, error)
	SearchImages(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error)
	SearchRawImages(ctx context.Context, q *v1.Query) ([]*storage.Image, error)

	CountImages(ctx context.Context) (int, error)
	GetImage(ctx context.Context, sha string) (*storage.Image, bool, error)
	GetImagesBatch(ctx context.Context, shas []string) ([]*storage.Image, error)

	UpsertImage(ctx context.Context, image *storage.Image) error

	DeleteImages(ctx context.Context, ids ...string) error
	Exists(ctx context.Context, id string) (bool, error)
}

func newDatastore(storage store.Store, bleveIndex bleve.Index, noUpdateTimestamps bool, risks riskDS.DataStore) (DataStore, error) {
	var searcher search.Searcher
	indexer := imageIndexer.New(bleveIndex)
	if features.Dackbox.Enabled() {
		searcher = search.New(storage,
			globalDackBox.GetGlobalDackBox(),
			cveIndexer.Singleton(),
			componentCVEEdgeIndexer.Singleton(),
			componentIndexer.Singleton(),
			imageComponentEdgeIndexer.Singleton(),
			imageIndexer.New(bleveIndex))
	} else {
		searcher = search.New(storage, nil, nil, nil, nil, nil, indexer)
	}

	ds, err := newDatastoreImpl(storage, indexer, searcher, risks)
	if err != nil {
		return nil, err
	}

	if err := ds.initializeRankers(); err != nil {
		return nil, errors.Wrap(err, "failed to initialize ranker")
	}

	return ds, nil
}

// NewBadger returns a new instance of DataStore using the input store, indexer, and searcher.
// noUpdateTimestamps controls whether timestamps are automatically updated when upserting images.
// This should be set to `false` except for some tests.
func NewBadger(db *badger.DB, bleveIndex bleve.Index, noUpdateTimestamps bool, risks riskDS.DataStore) (DataStore, error) {
	storage := badgerStore.New(db, noUpdateTimestamps)
	return newDatastore(storage, bleveIndex, noUpdateTimestamps, risks)
}
