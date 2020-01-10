package filtered

import (
	"context"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/blevesearch"
	"github.com/stackrox/rox/pkg/set"
)

// Filter represents a process of converting from one id-space to another.
type Filter interface {
	Apply(ctx context.Context, from ...string) ([]string, error)
}

// Searcher generates a Searcher with the input filter.
func Searcher(searcher blevesearch.UnsafeSearcher, filter Filter) search.Searcher {
	return search.Func(func(ctx context.Context, q *v1.Query) ([]search.Result, error) {
		results, err := searcher.Search(q)
		if err != nil {
			return results, err
		}

		filtered, err := filter.Apply(ctx, search.ResultsToIDs(results)...)
		if err != nil {
			return results, err
		}

		filteredResults := results[:0]
		filteredSet := set.NewStringSet(filtered...)
		for _, result := range results {
			if filteredSet.Contains(result.ID) {
				filteredResults = append(filteredResults, result)
			}
		}
		return filteredResults, nil
	})
}
