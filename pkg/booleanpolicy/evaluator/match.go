package evaluator

import (
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/traverseutil"
)

// A Match represents a single match.
// It contains the matched value, as well as the path
// within the object that was taken to reach the value.
type Match struct {
	Path  *traverseutil.Path
	Value interface{}
}

// A Result is the result of evaluating a query on an object.
type Result struct {
	Matches []Match
}

func mergeResults(results []*Result) *Result {
	if len(results) == 0 {
		return nil
	}
	var totalLen int
	for _, r := range results {
		totalLen += len(r.Matches)
	}
	merged := &Result{Matches: make([]Match, 0, totalLen)}
	for _, r := range results {
		merged.Matches = append(merged.Matches, r.Matches...)
	}
	return merged
}
