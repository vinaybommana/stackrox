// Code generated by blevebindings generator. DO NOT EDIT.

package index

import (
	mappings "github.com/stackrox/rox/central/clustercveedge/mappings"
	metrics "github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	batcher "github.com/stackrox/rox/pkg/batcher"
	blevehelper "github.com/stackrox/rox/pkg/blevehelper"
	ops "github.com/stackrox/rox/pkg/metrics"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
	"time"
)

const batchSize = 5000

const resourceName = "ClusterCVEEdge"

type indexerImpl struct {
	index *blevehelper.BleveWrapper
}

type clusterCVEEdgeWrapper struct {
	*storage.ClusterCVEEdge `json:"cluster_c_v_e_edge"`
	Type                    string `json:"type"`
}

func (b *indexerImpl) AddClusterCVEEdge(clustercveedge *storage.ClusterCVEEdge) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Add, "ClusterCVEEdge")
	if err := b.index.Index.Index(clustercveedge.GetId(), &clusterCVEEdgeWrapper{
		ClusterCVEEdge: clustercveedge,
		Type:           v1.SearchCategory_CLUSTER_VULN_EDGE.String(),
	}); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) AddClusterCVEEdges(clustercveedges []*storage.ClusterCVEEdge) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.AddMany, "ClusterCVEEdge")
	batchManager := batcher.New(len(clustercveedges), batchSize)
	for {
		start, end, ok := batchManager.Next()
		if !ok {
			break
		}
		if err := b.processBatch(clustercveedges[start:end]); err != nil {
			return err
		}
	}
	return nil
}

func (b *indexerImpl) processBatch(clustercveedges []*storage.ClusterCVEEdge) error {
	batch := b.index.NewBatch()
	for _, clustercveedge := range clustercveedges {
		if err := batch.Index(clustercveedge.GetId(), &clusterCVEEdgeWrapper{
			ClusterCVEEdge: clustercveedge,
			Type:           v1.SearchCategory_CLUSTER_VULN_EDGE.String(),
		}); err != nil {
			return err
		}
	}
	return b.index.Batch(batch)
}

func (b *indexerImpl) DeleteClusterCVEEdge(id string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Remove, "ClusterCVEEdge")
	if err := b.index.Delete(id); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) DeleteClusterCVEEdges(ids []string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.RemoveMany, "ClusterCVEEdge")
	batch := b.index.NewBatch()
	for _, id := range ids {
		batch.Delete(id)
	}
	if err := b.index.Batch(batch); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) ResetIndex() error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Reset, "ClusterCVEEdge")
	return blevesearch.ResetIndex(v1.SearchCategory_CLUSTER_VULN_EDGE, b.index.Index)
}

func (b *indexerImpl) Search(q *v1.Query, opts ...blevesearch.SearchOption) ([]search.Result, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Search, "ClusterCVEEdge")
	return blevesearch.RunSearchRequest(v1.SearchCategory_CLUSTER_VULN_EDGE, q, b.index.Index, mappings.OptionsMap, opts...)
}
