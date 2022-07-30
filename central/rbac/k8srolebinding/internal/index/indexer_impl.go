// Code generated by blevebindings generator. DO NOT EDIT.

package index

import (
	"bytes"
	"time"

	bleve "github.com/blevesearch/bleve"
	metrics "github.com/stackrox/rox/central/metrics"
	mappings "github.com/stackrox/rox/central/rbac/k8srolebinding/mappings"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/aux"
	storage "github.com/stackrox/rox/generated/storage"
	batcher "github.com/stackrox/rox/pkg/batcher"
	ops "github.com/stackrox/rox/pkg/metrics"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
)

const batchSize = 5000

const resourceName = "K8SRoleBinding"

type indexerImpl struct {
	index bleve.Index
}

type k8SRoleBindingWrapper struct {
	*storage.K8SRoleBinding `json:"k8s_role_binding"`
	Type                    string `json:"type"`
}

func (b *indexerImpl) AddK8SRoleBinding(k8srolebinding *storage.K8SRoleBinding) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Add, "K8SRoleBinding")
	if err := b.index.Index(k8srolebinding.GetId(), &k8SRoleBindingWrapper{
		K8SRoleBinding: k8srolebinding,
		Type:           v1.SearchCategory_ROLEBINDINGS.String(),
	}); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) AddK8SRoleBindings(k8srolebindings []*storage.K8SRoleBinding) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.AddMany, "K8SRoleBinding")
	batchManager := batcher.New(len(k8srolebindings), batchSize)
	for {
		start, end, ok := batchManager.Next()
		if !ok {
			break
		}
		if err := b.processBatch(k8srolebindings[start:end]); err != nil {
			return err
		}
	}
	return nil
}

func (b *indexerImpl) processBatch(k8srolebindings []*storage.K8SRoleBinding) error {
	batch := b.index.NewBatch()
	for _, k8srolebinding := range k8srolebindings {
		if err := batch.Index(k8srolebinding.GetId(), &k8SRoleBindingWrapper{
			K8SRoleBinding: k8srolebinding,
			Type:           v1.SearchCategory_ROLEBINDINGS.String(),
		}); err != nil {
			return err
		}
	}
	return b.index.Batch(batch)
}

func (b *indexerImpl) Count(q *auxpb.Query, opts ...blevesearch.SearchOption) (int, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Count, "K8SRoleBinding")
	return blevesearch.RunCountRequest(v1.SearchCategory_ROLEBINDINGS, q, b.index, mappings.OptionsMap, opts...)
}

func (b *indexerImpl) DeleteK8SRoleBinding(id string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Remove, "K8SRoleBinding")
	if err := b.index.Delete(id); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) DeleteK8SRoleBindings(ids []string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.RemoveMany, "K8SRoleBinding")
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
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Search, "K8SRoleBinding")
	return blevesearch.RunSearchRequest(v1.SearchCategory_ROLEBINDINGS, q, b.index, mappings.OptionsMap, opts...)
}
