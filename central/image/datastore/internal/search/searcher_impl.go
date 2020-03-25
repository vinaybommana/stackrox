package search

import (
	"context"

	componentCVEEdgeMappings "github.com/stackrox/rox/central/componentcveedge/mappings"
	"github.com/stackrox/rox/central/cve/cveedge"
	cveMappings "github.com/stackrox/rox/central/cve/mappings"
	"github.com/stackrox/rox/central/dackbox"
	pkgDeploymentSAC "github.com/stackrox/rox/central/deployment/sac"
	"github.com/stackrox/rox/central/image/datastore/internal/store"
	"github.com/stackrox/rox/central/image/index"
	pkgImageSAC "github.com/stackrox/rox/central/image/sac"
	componentMappings "github.com/stackrox/rox/central/imagecomponent/mappings"
	imageComponentEdgeMappings "github.com/stackrox/rox/central/imagecomponentedge/mappings"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/dackbox/graph"
	"github.com/stackrox/rox/pkg/derivedfields/counter"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/images/types"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/blevesearch"
	"github.com/stackrox/rox/pkg/search/compound"
	"github.com/stackrox/rox/pkg/search/derivedfields"
	"github.com/stackrox/rox/pkg/search/filtered"
	imageMappings "github.com/stackrox/rox/pkg/search/options/images"
	"github.com/stackrox/rox/pkg/search/paginated"
	"github.com/stackrox/rox/pkg/search/scoped"
	"github.com/stackrox/rox/pkg/search/sortfields"
)

var (
	defaultSortOption = &v1.QuerySortOption{
		Field: search.LastUpdatedTime.String(),
	}
	imagesSACSearchHelper = sac.ForResource(resources.Image).MustCreateSearchHelper(imageMappings.OptionsMap)

	componentOptionsMap = search.CombineOptionsMaps(componentMappings.OptionsMap).Remove(search.RiskScore)
	imageOnlyOptionsMap = search.Difference(
		imageMappings.OptionsMap,
		search.CombineOptionsMaps(
			imageComponentEdgeMappings.OptionsMap,
			componentOptionsMap,
			componentCVEEdgeMappings.OptionsMap,
			cveMappings.OptionsMap,
		),
	)
)

// searcherImpl provides an intermediary implementation layer for AlertStorage.
type searcherImpl struct {
	storage       store.Store
	indexer       index.Indexer
	graphProvider graph.Provider
	searcher      search.Searcher
}

// SearchImages retrieves SearchResults from the indexer and storage
func (ds *searcherImpl) SearchImages(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error) {
	images, results, err := ds.searchImages(ctx, q)
	if err != nil {
		return nil, err
	}
	protoResults := make([]*v1.SearchResult, 0, len(images))
	for i, image := range images {
		protoResults = append(protoResults, convertImage(image, results[i]))
	}
	return protoResults, nil
}

func (ds *searcherImpl) SearchListImages(ctx context.Context, q *v1.Query) ([]*storage.ListImage, error) {
	images, _, err := ds.searchImages(ctx, q)
	return images, err
}

// SearchRawImages retrieves SearchResults from the indexer and storage
func (ds *searcherImpl) SearchRawImages(ctx context.Context, q *v1.Query) ([]*storage.Image, error) {
	results, err := ds.Search(ctx, q)
	if err != nil {
		return nil, err
	}
	var images []*storage.Image
	for _, result := range results {
		image, exists, err := ds.storage.GetImage(result.ID)
		if err != nil {
			return nil, err
		}
		// The result may not exist if the object was deleted after the search
		if !exists {
			continue
		}
		images = append(images, image)
	}
	return images, nil
}

func (ds *searcherImpl) searchImages(ctx context.Context, q *v1.Query) ([]*storage.ListImage, []search.Result, error) {
	results, err := ds.Search(ctx, q)
	if err != nil {
		return nil, nil, err
	}
	var images []*storage.ListImage
	var newResults []search.Result
	for _, result := range results {
		image, exists, err := ds.storage.ListImage(result.ID)
		if err != nil {
			return nil, nil, err
		}
		// The result may not exist if the object was deleted after the search
		if !exists {
			continue
		}
		images = append(images, image)
		newResults = append(newResults, result)
	}
	return images, newResults, nil
}

func (ds *searcherImpl) Search(ctx context.Context, q *v1.Query) (res []search.Result, err error) {
	graph.Context(ctx, ds.graphProvider, func(inner context.Context) {
		res, err = ds.searcher.Search(inner, q)
	})
	return res, err
}

// ConvertImage returns proto search result from a image object and the internal search result
func convertImage(image *storage.ListImage, result search.Result) *v1.SearchResult {
	return &v1.SearchResult{
		Category:       v1.SearchCategory_IMAGES,
		Id:             types.NewDigest(image.GetId()).Digest(),
		Name:           image.GetName(),
		FieldToMatches: search.GetProtoMatchesMap(result.Matches),
		Score:          result.Score,
	}
}

// Format the search functionality of the indexer to be filtered (for sac) and paginated.
func formatSearcher(graphProvider graph.Provider,
	cveIndexer blevesearch.UnsafeSearcher,
	componentCVEEdgeIndexer blevesearch.UnsafeSearcher,
	componentIndexer blevesearch.UnsafeSearcher,
	imageComponentEdgeIndexer blevesearch.UnsafeSearcher,
	imageIndexer blevesearch.UnsafeSearcher) search.Searcher {
	var filteredSearcher search.Searcher
	if features.Dackbox.Enabled() {
		cveSearcher := blevesearch.WrapUnsafeSearcherAsSearcher(cveIndexer)
		componentCVEEdgeSearcher := blevesearch.WrapUnsafeSearcherAsSearcher(componentCVEEdgeIndexer)
		componentSearcher := blevesearch.WrapUnsafeSearcherAsSearcher(componentIndexer)
		imageComponentEdgeSearcher := blevesearch.WrapUnsafeSearcherAsSearcher(imageComponentEdgeIndexer)
		imageSearcher := blevesearch.WrapUnsafeSearcherAsSearcher(imageIndexer)

		compoundSearcher := getCompoundImageSearcher(
			cveSearcher,
			componentCVEEdgeSearcher,
			componentSearcher,
			imageComponentEdgeSearcher,
			imageSearcher)
		filteredSearcher = filtered.Searcher(cveedge.HandleCVEEdgeSearchQuery(compoundSearcher), pkgImageSAC.GetSACFilter(graphProvider))
	} else {
		filteredSearcher = imagesSACSearchHelper.FilteredSearcher(imageIndexer) // Make the UnsafeSearcher safe.
	}
	transformedSortSearcher := sortfields.TransformSortFields(filteredSearcher)
	derivedFieldSortedSearcher := wrapDerivedFieldSearcher(graphProvider, transformedSortSearcher)
	paginatedSearcher := paginated.Paginated(derivedFieldSortedSearcher)
	defaultSortedSearcher := paginated.WithDefaultSortOption(paginatedSearcher, defaultSortOption)
	return defaultSortedSearcher
}

func getCompoundImageSearcher(
	cveSearcher search.Searcher,
	componentCVEEdgeSearcher search.Searcher,
	componentSearcher search.Searcher,
	imageComponentEdgeSearcher search.Searcher,
	imageSearcher search.Searcher) search.Searcher {
	// The ordering of these is important, so do not change.
	return compound.NewSearcher([]compound.SearcherSpec{
		{
			Searcher:       scoped.WithScoping(cveSearcher, dackbox.ToCategory(v1.SearchCategory_VULNERABILITIES)),
			Transformation: dackbox.GraphTransformations[v1.SearchCategory_VULNERABILITIES][v1.SearchCategory_IMAGES],
			Options:        cveMappings.OptionsMap,
		},
		{
			Searcher:       scoped.WithScoping(componentCVEEdgeSearcher, dackbox.ToCategory(v1.SearchCategory_COMPONENT_VULN_EDGE)),
			Transformation: dackbox.GraphTransformations[v1.SearchCategory_COMPONENT_VULN_EDGE][v1.SearchCategory_IMAGES],
			Options:        componentCVEEdgeMappings.OptionsMap,
			LinkToPrev:     dackbox.GraphTransformations[v1.SearchCategory_VULNERABILITIES][v1.SearchCategory_COMPONENT_VULN_EDGE],
		},
		{
			Searcher:       scoped.WithScoping(componentSearcher, dackbox.ToCategory(v1.SearchCategory_IMAGE_COMPONENTS)),
			Transformation: dackbox.GraphTransformations[v1.SearchCategory_IMAGE_COMPONENTS][v1.SearchCategory_IMAGES],
			Options:        componentOptionsMap,
			LinkToPrev:     dackbox.GraphTransformations[v1.SearchCategory_COMPONENT_VULN_EDGE][v1.SearchCategory_IMAGE_COMPONENTS],
		},
		{
			Searcher:       scoped.WithScoping(imageComponentEdgeSearcher, dackbox.ToCategory(v1.SearchCategory_IMAGE_COMPONENT_EDGE)),
			Transformation: dackbox.GraphTransformations[v1.SearchCategory_IMAGE_COMPONENT_EDGE][v1.SearchCategory_IMAGES],
			Options:        imageComponentEdgeMappings.OptionsMap,
			LinkToPrev:     dackbox.GraphTransformations[v1.SearchCategory_IMAGE_COMPONENTS][v1.SearchCategory_IMAGE_COMPONENT_EDGE],
		},
		{
			IsDefault:  true,
			Searcher:   scoped.WithScoping(imageSearcher, dackbox.ToCategory(v1.SearchCategory_IMAGES)),
			Options:    imageOnlyOptionsMap,
			LinkToPrev: dackbox.GraphTransformations[v1.SearchCategory_IMAGE_COMPONENT_EDGE][v1.SearchCategory_IMAGES],
		},
	})
}

func wrapDerivedFieldSearcher(graphProvider graph.Provider, searcher search.Searcher) search.Searcher {
	if !features.Dackbox.Enabled() {
		return searcher
	}
	return derivedfields.CountSortedSearcher(searcher, map[string]counter.DerivedFieldCounter{
		search.DeploymentCount.String(): counter.NewGraphBasedDerivedFieldCounter(graphProvider, dackbox.ImageToDeploymentPath, pkgDeploymentSAC.GetSACFilter(graphProvider)),
	})
}
