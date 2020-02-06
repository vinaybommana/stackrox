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

// UnsafeSearcher generates a Searcher from an UnsafeSearcher by filtering its outputs with the input filter.
func UnsafeSearcher(searcher blevesearch.UnsafeSearcher, filters ...Filter) search.Searcher {
	return search.Func(func(ctx context.Context, q *v1.Query) ([]search.Result, error) {
		results, err := searcher.Search(q)
		if err != nil {
			return results, err
		}

		var filtered []string
		for _, filter := range filters {
			filtered, err = filter.Apply(ctx, search.ResultsToIDs(results)...)

			// If there is no error and if result is non-empty we assume that we evaluated on correct sac filter.
			// If we have receive error or if the length of result is 0, we try another filter.
			if err == nil && len(filtered) > 0 {
				break
			}
		}
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

// Searcher returns a new searcher based on the filtered output from the input Searcher.
func Searcher(searcher search.Searcher, filters ...Filter) search.Searcher {
	return search.Func(func(ctx context.Context, q *v1.Query) ([]search.Result, error) {
		var err error
		results, err := searcher.Search(ctx, q)
		if err != nil {
			return results, err
		}

		var filtered []string
		for _, filter := range filters {
			filtered, err = filter.Apply(ctx, search.ResultsToIDs(results)...)

			// If there is no error and if result is non-empty we assume that we evaluated on correct sac filter.
			// If we have receive error or if the length of result is 0, we try another filter.
			if err == nil && len(filtered) > 0 {
				break
			}
		}
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
