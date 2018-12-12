package datastore

import (
	"github.com/stackrox/rox/central/policy/index"
	"github.com/stackrox/rox/central/policy/search"
	"github.com/stackrox/rox/central/policy/store"
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
)

// DataStore is an intermediary to PolicyStorage.
//go:generate mockgen-wrapper DataStore
type DataStore interface {
	SearchPolicies(q *v1.Query) ([]*v1.SearchResult, error)
	SearchRawPolicies(q *v1.Query) ([]*storage.Policy, error)

	GetPolicy(id string) (*storage.Policy, bool, error)
	GetPolicies() ([]*storage.Policy, error)

	AddPolicy(*storage.Policy) (string, error)
	UpdatePolicy(*storage.Policy) error
	RemovePolicy(id string) error
	RenamePolicyCategory(request *v1.RenamePolicyCategoryRequest) error
	DeletePolicyCategory(request *v1.DeletePolicyCategoryRequest) error
}

// New returns a new instance of DataStore using the input store, indexer, and searcher.
func New(storage store.Store, indexer index.Indexer, searcher search.Searcher) DataStore {
	return &datastoreImpl{
		storage:  storage,
		indexer:  indexer,
		searcher: searcher,
	}
}
