// Code generated by blevebindings generator. DO NOT EDIT.

package index

import (
	"bytes"
	"time"

	bleve "github.com/blevesearch/bleve"
	metrics "github.com/stackrox/rox/central/metrics"
	mappings "github.com/stackrox/rox/central/policy/index/mappings"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/aux"
	storage "github.com/stackrox/rox/generated/storage"
	batcher "github.com/stackrox/rox/pkg/batcher"
	ops "github.com/stackrox/rox/pkg/metrics"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
)

const batchSize = 5000

const resourceName = "Policy"

type indexerImpl struct {
	index bleve.Index
}

type policyWrapper struct {
	*storage.Policy `json:"policy"`
	Type            string `json:"type"`
}

func (b *indexerImpl) AddPolicy(policy *storage.Policy) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Add, "Policy")
	if err := b.index.Index(policy.GetId(), &policyWrapper{
		Policy: policy,
		Type:   v1.SearchCategory_POLICIES.String(),
	}); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) AddPolicies(policies []*storage.Policy) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.AddMany, "Policy")
	batchManager := batcher.New(len(policies), batchSize)
	for {
		start, end, ok := batchManager.Next()
		if !ok {
			break
		}
		if err := b.processBatch(policies[start:end]); err != nil {
			return err
		}
	}
	return nil
}

func (b *indexerImpl) processBatch(policies []*storage.Policy) error {
	batch := b.index.NewBatch()
	for _, policy := range policies {
		if err := batch.Index(policy.GetId(), &policyWrapper{
			Policy: policy,
			Type:   v1.SearchCategory_POLICIES.String(),
		}); err != nil {
			return err
		}
	}
	return b.index.Batch(batch)
}

func (b *indexerImpl) Count(q *auxpb.Query, opts ...blevesearch.SearchOption) (int, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Count, "Policy")
	return blevesearch.RunCountRequest(v1.SearchCategory_POLICIES, q, b.index, mappings.OptionsMap, opts...)
}

func (b *indexerImpl) DeletePolicy(id string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Remove, "Policy")
	if err := b.index.Delete(id); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) DeletePolicies(ids []string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.RemoveMany, "Policy")
	batch := b.index.NewBatch()
	for _, id := range ids {
		batch.Delete(id)
	}
	if err := b.index.Batch(batch); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) MarkInitialIndexingComplete() error {
	return b.index.SetInternal([]byte(resourceName), []byte("old"))
}

func (b *indexerImpl) NeedsInitialIndexing() (bool, error) {
	data, err := b.index.GetInternal([]byte(resourceName))
	if err != nil {
		return false, err
	}
	return !bytes.Equal([]byte("old"), data), nil
}

func (b *indexerImpl) Search(q *auxpb.Query, opts ...blevesearch.SearchOption) ([]search.Result, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Search, "Policy")
	return blevesearch.RunSearchRequest(v1.SearchCategory_POLICIES, q, b.index, mappings.OptionsMap, opts...)
}
