package datastore

import (
	"context"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/cve/index"
	"github.com/stackrox/rox/central/cve/search"
	"github.com/stackrox/rox/central/cve/store"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/sac"
	searchPkg "github.com/stackrox/rox/pkg/search"
)

var (
	// TODO: Need to setup sac for CVEs correctly instead of relying on global access.
	imagesSAC = sac.ForResource(resources.Image)
)

type datastoreImpl struct {
	storage  store.Store
	indexer  index.Indexer
	searcher search.Searcher
}

func (ds *datastoreImpl) Search(ctx context.Context, q *v1.Query) ([]searchPkg.Result, error) {
	if ok, err := imagesSAC.ReadAllowed(ctx); !ok || err != nil {
		return nil, err
	}
	return ds.searcher.Search(ctx, q)
}

func (ds *datastoreImpl) SearchCVEs(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error) {
	if ok, err := imagesSAC.ReadAllowed(ctx); !ok || err != nil {
		return nil, err
	}
	return ds.searcher.SearchCVEs(ctx, q)
}

func (ds *datastoreImpl) SearchRawCVEs(ctx context.Context, q *v1.Query) ([]*storage.CVE, error) {
	if ok, err := imagesSAC.ReadAllowed(ctx); !ok || err != nil {
		return nil, err
	}
	cves, err := ds.searcher.SearchRawCVEs(ctx, q)
	if err != nil {
		return nil, err
	}
	return cves, nil
}

func (ds *datastoreImpl) Count(ctx context.Context) (int, error) {
	if ok, err := imagesSAC.ReadAllowed(ctx); err != nil || !ok {
		return 0, err
	}
	return ds.storage.Count()
}

func (ds *datastoreImpl) Get(ctx context.Context, id string) (*storage.CVE, bool, error) {
	if ok, err := imagesSAC.ReadAllowed(ctx); err != nil || !ok {
		return nil, false, err
	}
	cve, found, err := ds.storage.Get(id)
	if err != nil || !found {
		return nil, false, err
	}
	return cve, true, nil
}

func (ds *datastoreImpl) Exists(ctx context.Context, id string) (bool, error) {
	if ok, err := imagesSAC.ReadAllowed(ctx); err != nil || !ok {
		return false, err
	}
	found, err := ds.storage.Exists(id)
	if err != nil || !found {
		return false, err
	}
	return true, nil
}

func (ds *datastoreImpl) GetBatch(ctx context.Context, ids []string) ([]*storage.CVE, error) {
	if ok, err := imagesSAC.ReadAllowed(ctx); err != nil || !ok {
		return nil, err
	}
	cves, _, err := ds.storage.GetBatch(ids)
	if err != nil {
		return nil, err
	}
	return cves, nil
}

func (ds *datastoreImpl) Upsert(ctx context.Context, cves ...*storage.CVE) error {
	if len(cves) == 0 {
		return nil
	}
	if ok, err := imagesSAC.WriteAllowed(ctx); err != nil {
		return err
	} else if !ok {
		return errors.New("permission denied")
	}

	// Load the supressed value for any CVEs already present.
	ids := make([]string, 0, len(cves))
	for _, cve := range cves {
		ids = append(ids, cve.GetId())
	}
	currentCVEs, _, err := ds.storage.GetBatch(ids)
	if err != nil {
		return err
	}
	var currentIndex int
	for newIndex := 0; newIndex < len(cves) && currentIndex < len(currentCVEs); newIndex++ {
		if currentCVEs[currentIndex].GetId() == cves[newIndex].GetId() {
			cves[newIndex].Supressed = currentCVEs[currentIndex].Supressed
			currentIndex++
		}
	}

	// Store the new CVE data.
	if err := ds.storage.Upsert(cves...); err != nil {
		return err
	}
	return ds.indexer.AddCVEs(cves)
}

func (ds *datastoreImpl) Supress(ctx context.Context, ids ...string) error {
	if ok, err := imagesSAC.WriteAllowed(ctx); err != nil {
		return err
	} else if !ok {
		return errors.New("permission denied")
	}

	cves, _, err := ds.storage.GetBatch(ids)
	if err != nil {
		return err
	}
	for _, cve := range cves {
		cve.Supressed = true
	}
	if err := ds.storage.Upsert(cves...); err != nil {
		return err
	}
	return ds.indexer.AddCVEs(cves)
}

func (ds *datastoreImpl) Unsupress(ctx context.Context, ids ...string) error {
	if ok, err := imagesSAC.WriteAllowed(ctx); err != nil {
		return err
	} else if !ok {
		return errors.New("permission denied")
	}

	cves, _, err := ds.storage.GetBatch(ids)
	if err != nil {
		return err
	}
	for _, cve := range cves {
		cve.Supressed = false
	}
	if err := ds.storage.Upsert(cves...); err != nil {
		return err
	}
	return ds.indexer.AddCVEs(cves)
}

func (ds *datastoreImpl) Delete(ctx context.Context, ids ...string) error {
	if ok, err := imagesSAC.WriteAllowed(ctx); err != nil {
		return err
	} else if !ok {
		return errors.New("permission denied")
	}

	if err := ds.storage.Delete(ids...); err != nil {
		return err
	}
	return ds.indexer.DeleteCVEs(ids)
}
