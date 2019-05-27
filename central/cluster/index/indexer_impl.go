// Code generated by blevebindings generator. DO NOT EDIT.

package index

import (
	bleve "github.com/blevesearch/bleve"
	mappings "github.com/stackrox/rox/central/cluster/index/mappings"
	metrics "github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	batcher "github.com/stackrox/rox/pkg/batcher"
	ops "github.com/stackrox/rox/pkg/metrics"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
	"time"
)

const batchSize = 5000

type indexerImpl struct {
	index bleve.Index
}

type clusterWrapper struct {
	*storage.Cluster `json:"cluster"`
	Type             string `json:"type"`
}

func (b *indexerImpl) AddCluster(cluster *storage.Cluster) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Add, "Cluster")
	return b.index.Index(cluster.GetId(), &clusterWrapper{
		Cluster: cluster,
		Type:    v1.SearchCategory_CLUSTERS.String(),
	})
}

func (b *indexerImpl) AddClusters(clusters []*storage.Cluster) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.AddMany, "Cluster")
	batchManager := batcher.New(len(clusters), batchSize)
	for {
		start, end, ok := batchManager.Next()
		if !ok {
			break
		}
		if err := b.processBatch(clusters[start:end]); err != nil {
			return err
		}
	}
	return nil
}

func (b *indexerImpl) processBatch(clusters []*storage.Cluster) error {
	batch := b.index.NewBatch()
	for _, cluster := range clusters {
		if err := batch.Index(cluster.GetId(), &clusterWrapper{
			Cluster: cluster,
			Type:    v1.SearchCategory_CLUSTERS.String(),
		}); err != nil {
			return err
		}
	}
	return b.index.Batch(batch)
}

func (b *indexerImpl) DeleteCluster(id string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Remove, "Cluster")
	return b.index.Delete(id)
}

func (b *indexerImpl) Search(q *v1.Query) ([]search.Result, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Search, "Cluster")
	return blevesearch.RunSearchRequest(v1.SearchCategory_CLUSTERS, q, b.index, mappings.OptionsMap)
}
