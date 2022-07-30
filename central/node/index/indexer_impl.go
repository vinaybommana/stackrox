// Code generated by blevebindings generator. DO NOT EDIT.

package index

import (
	"bytes"
	"time"

	bleve "github.com/blevesearch/bleve"
	metrics "github.com/stackrox/rox/central/metrics"
	mappings "github.com/stackrox/rox/central/node/index/mappings"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/aux"
	storage "github.com/stackrox/rox/generated/storage"
	batcher "github.com/stackrox/rox/pkg/batcher"
	ops "github.com/stackrox/rox/pkg/metrics"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
)

const batchSize = 5000

const resourceName = "Node"

type indexerImpl struct {
	index bleve.Index
}

type nodeWrapper struct {
	*storage.Node `json:"node"`
	Type          string `json:"type"`
}

func (b *indexerImpl) AddNode(node *storage.Node) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Add, "Node")
	if err := b.index.Index(node.GetId(), &nodeWrapper{
		Node: node,
		Type: v1.SearchCategory_NODES.String(),
	}); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) AddNodes(nodes []*storage.Node) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.AddMany, "Node")
	batchManager := batcher.New(len(nodes), batchSize)
	for {
		start, end, ok := batchManager.Next()
		if !ok {
			break
		}
		if err := b.processBatch(nodes[start:end]); err != nil {
			return err
		}
	}
	return nil
}

func (b *indexerImpl) processBatch(nodes []*storage.Node) error {
	batch := b.index.NewBatch()
	for _, node := range nodes {
		if err := batch.Index(node.GetId(), &nodeWrapper{
			Node: node,
			Type: v1.SearchCategory_NODES.String(),
		}); err != nil {
			return err
		}
	}
	return b.index.Batch(batch)
}

func (b *indexerImpl) Count(q *auxpb.Query, opts ...blevesearch.SearchOption) (int, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Count, "Node")
	return blevesearch.RunCountRequest(v1.SearchCategory_NODES, q, b.index, mappings.OptionsMap, opts...)
}

func (b *indexerImpl) DeleteNode(id string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Remove, "Node")
	if err := b.index.Delete(id); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) DeleteNodes(ids []string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.RemoveMany, "Node")
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
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Search, "Node")
	return blevesearch.RunSearchRequest(v1.SearchCategory_NODES, q, b.index, mappings.OptionsMap, opts...)
}
