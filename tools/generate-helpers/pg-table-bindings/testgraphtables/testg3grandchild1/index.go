// Code generated by pg-bindings generator. DO NOT EDIT.
package postgres

import (
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	metrics "github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/aux"
	storage "github.com/stackrox/rox/generated/storage"
	ops "github.com/stackrox/rox/pkg/metrics"
	search "github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/blevesearch"
	"github.com/stackrox/rox/pkg/search/postgres"
	"github.com/stackrox/rox/pkg/search/postgres/mapping"
)

func init() {
	mapping.RegisterCategoryToTable(v1.SearchCategory(67), schema)
}

// NewIndexer returns new indexer for `storage.TestG3GrandChild1`.
func NewIndexer(db *pgxpool.Pool) *indexerImpl {
	return &indexerImpl{
		db: db,
	}
}

type indexerImpl struct {
	db *pgxpool.Pool
}

func (b *indexerImpl) Count(q *auxpb.Query, opts ...blevesearch.SearchOption) (int, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Count, "TestG3GrandChild1")

	return postgres.RunCountRequest(v1.SearchCategory(67), q, b.db)
}

func (b *indexerImpl) Search(q *auxpb.Query, opts ...blevesearch.SearchOption) ([]search.Result, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Search, "TestG3GrandChild1")

	return postgres.RunSearchRequest(v1.SearchCategory(67), q, b.db)
}

//// Stubs for satisfying interfaces

func (b *indexerImpl) AddTestG3GrandChild1(deployment *storage.TestG3GrandChild1) error {
	return nil
}

func (b *indexerImpl) AddTestG3GrandChild1s(_ []*storage.TestG3GrandChild1) error {
	return nil
}

func (b *indexerImpl) DeleteTestG3GrandChild1(id string) error {
	return nil
}

func (b *indexerImpl) DeleteTestG3GrandChild1s(_ []string) error {
	return nil
}

func (b *indexerImpl) MarkInitialIndexingComplete() error {
	return nil
}

func (b *indexerImpl) NeedsInitialIndexing() (bool, error) {
	return false, nil
}
