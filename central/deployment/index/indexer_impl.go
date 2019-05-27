// Code generated by blevebindings generator. DO NOT EDIT.

package index

import (
	bleve "github.com/blevesearch/bleve"
	mappings "github.com/stackrox/rox/central/deployment/mappings"
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

type deploymentWrapper struct {
	*storage.Deployment `json:"deployment"`
	Type                string `json:"type"`
}

func (b *indexerImpl) AddDeployment(deployment *storage.Deployment) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Add, "Deployment")
	return b.index.Index(deployment.GetId(), &deploymentWrapper{
		Deployment: deployment,
		Type:       v1.SearchCategory_DEPLOYMENTS.String(),
	})
}

func (b *indexerImpl) AddDeployments(deployments []*storage.Deployment) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.AddMany, "Deployment")
	batchManager := batcher.New(len(deployments), batchSize)
	for {
		start, end, ok := batchManager.Next()
		if !ok {
			break
		}
		if err := b.processBatch(deployments[start:end]); err != nil {
			return err
		}
	}
	return nil
}

func (b *indexerImpl) processBatch(deployments []*storage.Deployment) error {
	batch := b.index.NewBatch()
	for _, deployment := range deployments {
		if err := batch.Index(deployment.GetId(), &deploymentWrapper{
			Deployment: deployment,
			Type:       v1.SearchCategory_DEPLOYMENTS.String(),
		}); err != nil {
			return err
		}
	}
	return b.index.Batch(batch)
}

func (b *indexerImpl) DeleteDeployment(id string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Remove, "Deployment")
	return b.index.Delete(id)
}

func (b *indexerImpl) Search(q *v1.Query) ([]search.Result, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Search, "Deployment")
	return blevesearch.RunSearchRequest(v1.SearchCategory_DEPLOYMENTS, q, b.index, mappings.OptionsMap)
}
