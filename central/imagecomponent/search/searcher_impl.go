package search

import (
	"context"

	componentCVEEdgeMappings "github.com/stackrox/rox/central/componentcveedge/mappings"
	pkgComponentCVEEdgeSAC "github.com/stackrox/rox/central/componentcveedge/sac"
	cveDackBox "github.com/stackrox/rox/central/cve/dackbox"
	cveMappings "github.com/stackrox/rox/central/cve/mappings"
	pkgCVESAC "github.com/stackrox/rox/central/cve/sac"
	imageDackBox "github.com/stackrox/rox/central/image/dackbox"
	imageMappings "github.com/stackrox/rox/central/image/mappings"
	pkgImageSAC "github.com/stackrox/rox/central/image/sac"
	componentDackBox "github.com/stackrox/rox/central/imagecomponent/dackbox"
	"github.com/stackrox/rox/central/imagecomponent/index"
	componentMappings "github.com/stackrox/rox/central/imagecomponent/mappings"
	pkgComponentSAC "github.com/stackrox/rox/central/imagecomponent/sac"
	"github.com/stackrox/rox/central/imagecomponent/store"
	imageComponentEdgeMappings "github.com/stackrox/rox/central/imagecomponentedge/mappings"
	pkgImageComponentEdgeSAC "github.com/stackrox/rox/central/imagecomponentedge/sac"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/blevesearch"
	"github.com/stackrox/rox/pkg/search/compound"
	"github.com/stackrox/rox/pkg/search/filtered"
	"github.com/stackrox/rox/pkg/search/idspace"
	"github.com/stackrox/rox/pkg/search/paginated"
)

var (
	defaultSortOption = &v1.QuerySortOption{
		Field: search.ComponentName.String(),
	}
)

type searcherImpl struct {
	storage  store.Store
	indexer  index.Indexer
	searcher search.Searcher
}

func (ds *searcherImpl) Search(ctx context.Context, q *v1.Query) ([]search.Result, error) {
	return ds.getSearchResults(ctx, q)
}

func (ds *searcherImpl) SearchImageComponents(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error) {
	results, err := ds.getSearchResults(ctx, q)
	if err != nil {
		return nil, err
	}
	return ds.resultsToSearchResults(results)
}

func (ds *searcherImpl) SearchRawImageComponents(ctx context.Context, q *v1.Query) ([]*storage.ImageComponent, error) {
	return ds.searchImageComponents(ctx, q)
}

func (ds *searcherImpl) getSearchResults(ctx context.Context, q *v1.Query) ([]search.Result, error) {
	return ds.searcher.Search(ctx, q)
}

func (ds *searcherImpl) resultsToImageComponents(results []search.Result) ([]*storage.ImageComponent, []int, error) {
	return ds.storage.GetBatch(search.ResultsToIDs(results))
}

func (ds *searcherImpl) resultsToSearchResults(results []search.Result) ([]*v1.SearchResult, error) {
	components, missingIndices, err := ds.resultsToImageComponents(results)
	if err != nil {
		return nil, err
	}
	results = search.RemoveMissingResults(results, missingIndices)
	return convertMany(components, results), nil
}

func convertMany(components []*storage.ImageComponent, results []search.Result) []*v1.SearchResult {
	outputResults := make([]*v1.SearchResult, len(components))
	for index, sar := range components {
		outputResults[index] = convertOne(sar, &results[index])
	}
	return outputResults
}

func convertOne(component *storage.ImageComponent, result *search.Result) *v1.SearchResult {
	return &v1.SearchResult{
		Category:       v1.SearchCategory_IMAGE_COMPONENTS,
		Id:             component.GetId(),
		Name:           component.GetName(),
		FieldToMatches: search.GetProtoMatchesMap(result.Matches),
		Score:          result.Score,
	}
}

// Format the search functionality of the indexer to be filtered (for sac) and paginated.
func formatSearcher(graphProvider idspace.GraphProvider,
	cveIndexer blevesearch.UnsafeSearcher,
	componentCVEEdgeIndexer blevesearch.UnsafeSearcher,
	componentIndexer blevesearch.UnsafeSearcher,
	imageComponentEdgeIndexer blevesearch.UnsafeSearcher,
	imageIndexer blevesearch.UnsafeSearcher) search.Searcher {
	cveSearcher := filtered.Searcher(cveIndexer, pkgCVESAC.CVESACFilter)
	componentCVEEdgeSearcher := filtered.Searcher(componentCVEEdgeIndexer, pkgComponentCVEEdgeSAC.ComponentCVEEdgeSACFilter)
	componentSearcher := filtered.Searcher(componentIndexer, pkgComponentSAC.ImageComponentSACFilter)
	imageComponentEdgeSearcher := filtered.Searcher(imageComponentEdgeIndexer, pkgImageComponentEdgeSAC.ImageComponentEdgeSACFilter)
	imageSearcher := filtered.Searcher(imageIndexer, pkgImageSAC.ImageSACFilter)

	compoundSearcher := getCompoundComponentSearcher(graphProvider,
		cveSearcher,
		componentCVEEdgeSearcher,
		componentSearcher,
		imageComponentEdgeSearcher,
		imageSearcher)
	paginatedSearcher := paginated.Paginated(compoundSearcher)
	defaultSortedSearcher := paginated.WithDefaultSortOption(paginatedSearcher, defaultSortOption)
	return defaultSortedSearcher
}

func (ds *searcherImpl) searchImageComponents(ctx context.Context, q *v1.Query) ([]*storage.ImageComponent, error) {
	results, err := ds.Search(ctx, q)
	if err != nil {
		return nil, err
	}

	ids := search.ResultsToIDs(results)
	components, _, err := ds.storage.GetBatch(ids)
	if err != nil {
		return nil, err
	}
	return components, nil
}

func getCompoundComponentSearcher(graphProvider idspace.GraphProvider,
	cveSearcher search.Searcher,
	componentCVEEdgeSearcher search.Searcher,
	componentSearcher search.Searcher,
	imageComponentEdgeSearcher search.Searcher,
	imageSearcher search.Searcher) search.Searcher {

	return compound.NewSearcher([]compound.SearcherSpec{
		{
			Searcher: idspace.TransformIDs(cveSearcher, idspace.NewBackwardGraphTransformer(graphProvider,
				[][]byte{cveDackBox.Bucket,
					componentDackBox.Bucket,
				})),
			Options: cveMappings.OptionsMap,
		},
		{
			Searcher: idspace.TransformIDs(componentCVEEdgeSearcher, idspace.NewEdgeToParentTransformer()),
			Options:  componentCVEEdgeMappings.OptionsMap,
		},
		{
			Searcher: componentSearcher,
			Options:  componentMappings.OptionsMap,
		},
		{
			Searcher: idspace.TransformIDs(imageComponentEdgeSearcher, idspace.NewEdgeToChildTransformer()),
			Options:  imageComponentEdgeMappings.OptionsMap,
		},
		{
			Searcher: idspace.TransformIDs(imageSearcher, idspace.NewForwardGraphTransformer(graphProvider,
				[][]byte{imageDackBox.Bucket,
					componentDackBox.Bucket,
				})),
			Options: imageMappings.OptionsMap,
		},
	}...)
}
