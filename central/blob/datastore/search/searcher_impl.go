package search

import (
	"context"

	"github.com/stackrox/rox/central/blob/datastore/index"
	"github.com/stackrox/rox/central/blob/datastore/store"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/blevesearch"
	"github.com/stackrox/rox/pkg/search/paginated"
)

type searcherImpl struct {
	storage           store.Store
	indexer           index.Indexer
	formattedSearcher search.Searcher
}

func (s *searcherImpl) SearchIDs(ctx context.Context, q *v1.Query) ([]string, error) {
	var (
		results []search.Result
		err     error
	)
	results, err = s.indexer.Search(ctx, q)
	if err != nil || len(results) == 0 {
		return nil, err
	}
	ids := search.ResultsToIDs(results)
	return ids, nil
}

func (s *searcherImpl) SearchBlobsWithoutData(ctx context.Context, q *v1.Query) ([]*storage.Blob, error) {
	ids, err := s.SearchIDs(ctx, q)
	if err != nil {
		return nil, err
	}
	blobs, _, err := s.storage.GetManyBlobMetadata(ctx, ids)
	if err != nil {
		return nil, err
	}
	return blobs, nil
}

func (s *searcherImpl) Search(ctx context.Context, q *v1.Query) ([]search.Result, error) {
	return s.formattedSearcher.Search(ctx, q)
}

// Count returns the number of search results from the query
func (s *searcherImpl) Count(ctx context.Context, q *v1.Query) (int, error) {
	return s.formattedSearcher.Count(ctx, q)
}

// Helper functions which format our searching.
///////////////////////////////////////////////

func formatSearcher(unsafeSearcher blevesearch.UnsafeSearcher) search.Searcher {
	searcher := blevesearch.WrapUnsafeSearcherAsSearcher(unsafeSearcher)
	paginatedSearcher := paginated.Paginated(searcher)
	return paginatedSearcher
}
