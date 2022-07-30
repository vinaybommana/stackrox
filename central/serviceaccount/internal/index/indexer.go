// Code generated by blevebindings generator. DO NOT EDIT.

package index

import (
	bleve "github.com/blevesearch/bleve"
	"github.com/stackrox/rox/generated/aux"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
)

type Indexer interface {
	AddServiceAccount(serviceaccount *storage.ServiceAccount) error
	AddServiceAccounts(serviceaccounts []*storage.ServiceAccount) error
	Count(q *auxpb.Query, opts ...blevesearch.SearchOption) (int, error)
	DeleteServiceAccount(id string) error
	DeleteServiceAccounts(ids []string) error
	MarkInitialIndexingComplete() error
	NeedsInitialIndexing() (bool, error)
	Search(q *auxpb.Query, opts ...blevesearch.SearchOption) ([]search.Result, error)
}

func New(index bleve.Index) Indexer {
	return &indexerImpl{index: index}
}
