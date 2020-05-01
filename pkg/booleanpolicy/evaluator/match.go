package evaluator

import (
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/traverseutil"
)

// A Match represents a single match.
// It contains the matched value, as well as the path
// within the object that was taken to reach the value.
type Match struct {
	Path *traverseutil.Path
	// Values represents a list of human-friendly representations of the matched value.
	Values []string
}

// GetPath implements the traverseutil.PathHolder interface.
func (m Match) GetPath() *traverseutil.Path {
	return m.Path
}

// A Result is the result of evaluating a query on an object.
type Result struct {
	Matches map[string][]Match
}

func newResult() *Result {
	return &Result{Matches: make(map[string][]Match)}
}

func mergeResults(results []*Result) *Result {
	if len(results) == 0 {
		return nil
	}

	merged := &Result{Matches: make(map[string][]Match)}
	for _, r := range results {
		for field, matches := range r.Matches {
			merged.Matches[field] = append(merged.Matches[field], matches...)
		}
	}
	return merged
}
