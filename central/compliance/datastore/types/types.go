package types

import (
	"github.com/stackrox/rox/generated/storage"
)

// GetFlags controls the behavior of the Get... methods of a Store.
type GetFlags int32

const (
	// WithMessageStrings will cause compliance results to be loaded with message strings.
	WithMessageStrings GetFlags = 1 << iota
	// RequireMessageStrings implies WithMessageStrings, and additionally fails with an error if any message strings
	// could not be loaded.
	RequireMessageStrings
)

// ResultsWithStatus returns the last successful results, as well as the metadata for the recent (i.e., since the
// last successful results) failed results.
type ResultsWithStatus struct {
	LastSuccessfulResults *storage.ComplianceRunResults
	FailedRuns            []*storage.ComplianceRunMetadata
}
