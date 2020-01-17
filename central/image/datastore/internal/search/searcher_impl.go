package search

import (
	"context"

	componentCVEEdgeMappings "github.com/stackrox/rox/central/componentcveedge/mappings"
	pkgComponentCVESAC "github.com/stackrox/rox/central/componentcveedge/sac"
	cveDackBox "github.com/stackrox/rox/central/cve/dackbox"
	cveMappings "github.com/stackrox/rox/central/cve/mappings"
	pkgCVESAC "github.com/stackrox/rox/central/cve/sac"
	imageDackBox "github.com/stackrox/rox/central/image/dackbox"
	"github.com/stackrox/rox/central/image/datastore/internal/store"
	"github.com/stackrox/rox/central/image/index"
	pkgImageSAC "github.com/stackrox/rox/central/image/sac"
	componentDackBox "github.com/stackrox/rox/central/imagecomponent/dackbox"
	componentMappings "github.com/stackrox/rox/central/imagecomponent/mappings"
	pkgComponentSAC "github.com/stackrox/rox/central/imagecomponent/sac"
	imageComponentEdgeMappings "github.com/stackrox/rox/central/imagecomponentedge/mappings"
	pkgImageComponentEdgeSAC "github.com/stackrox/rox/central/imagecomponentedge/sac"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/images/types"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/blevesearch"
	"github.com/stackrox/rox/pkg/search/compound"
	"github.com/stackrox/rox/pkg/search/filtered"
	"github.com/stackrox/rox/pkg/search/idspace"
	mappings "github.com/stackrox/rox/pkg/search/options/images"
	"github.com/stackrox/rox/pkg/search/paginated"
	"github.com/stackrox/rox/pkg/search/sortfields"
)

var (
	defaultSortOption = &v1.QuerySortOption{
		Field: search.LastUpdatedTime.String(),
	}
	imagesSACSearchHelper = sac.ForResource(resources.Image).MustCreateSearchHelper(mappings.OptionsMap)
)

// searcherImpl provides an intermediary implementation layer for AlertStorage.
type searcherImpl struct {
	storage  store.Store
	indexer  index.Indexer
	searcher search.Searcher
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

func (ds *searcherImpl) Search(ctx context.Context, q *v1.Query) ([]search.Result, error) {
	return ds.searcher.Search(ctx, q)
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
func formatSearcher(graphProvider idspace.GraphProvider,
	cveIndexer blevesearch.UnsafeSearcher,
	componentCVEEdgeIndexer blevesearch.UnsafeSearcher,
	componentIndexer blevesearch.UnsafeSearcher,
	imageComponentEdgeIndexer blevesearch.UnsafeSearcher,
	imageIndexer blevesearch.UnsafeSearcher) search.Searcher {
	var paginatedSearcher search.Searcher
	if features.Dackbox.Enabled() {
		cveSearcher := filtered.Searcher(cveIndexer, pkgCVESAC.CVESACFilter)
		componentCVEEdgeSearcher := filtered.Searcher(componentCVEEdgeIndexer, pkgComponentCVESAC.ComponentCVEEdgeSACFilter)
		componentSearcher := filtered.Searcher(componentIndexer, pkgComponentSAC.ImageComponentSACFilter)
		imageComponentEdgeSearcher := filtered.Searcher(imageComponentEdgeIndexer, pkgImageComponentEdgeSAC.ImageComponentEdgeSACFilter)
		imageSearcher := filtered.Searcher(imageIndexer, pkgImageSAC.ImageSACFilter)

		compoundSearcher := getCompoundImageSearcher(graphProvider,
			cveSearcher,
			componentCVEEdgeSearcher,
			componentSearcher,
			imageComponentEdgeSearcher,
			imageSearcher)
		paginatedSearcher = paginated.Paginated(compoundSearcher)
	} else {
		filteredSearcher := imagesSACSearchHelper.FilteredSearcher(imageIndexer) // Make the UnsafeSearcher safe.
		transformedSortSearcher := sortfields.TransformSortFields(filteredSearcher)
		paginatedSearcher = paginated.Paginated(transformedSortSearcher)
	}
	defaultSortedSearcher := paginated.WithDefaultSortOption(paginatedSearcher, defaultSortOption)
	return defaultSortedSearcher
}

func getCompoundImageSearcher(graphProvider idspace.GraphProvider,
	cveSearcher search.Searcher,
	componentCVEEdgeSearcher search.Searcher,
	componentSearcher search.Searcher,
	imageComponentEdgeSearcher search.Searcher,
	imageSearcher search.Searcher) search.Searcher {
	componentCVEEdgeToComponentSearcher := idspace.TransformIDs(componentCVEEdgeSearcher, idspace.NewEdgeToParentTransformer())
	return compound.NewSearcher([]compound.SearcherSpec{
		{
			Searcher: idspace.TransformIDs(cveSearcher, idspace.NewBackwardGraphTransformer(graphProvider,
				[][]byte{cveDackBox.Bucket,
					componentDackBox.Bucket,
					imageDackBox.Bucket,
				})),
			Options: cveMappings.OptionsMap,
		},
		{
			Searcher: idspace.TransformIDs(componentCVEEdgeToComponentSearcher, idspace.NewBackwardGraphTransformer(graphProvider,
				[][]byte{componentDackBox.Bucket,
					imageDackBox.Bucket,
				})),
			Options: componentCVEEdgeMappings.OptionsMap,
		},
		{
			Searcher: idspace.TransformIDs(componentSearcher, idspace.NewBackwardGraphTransformer(graphProvider,
				[][]byte{componentDackBox.Bucket,
					imageDackBox.Bucket,
				})),
			Options: componentMappings.OptionsMap,
		},
		{
			Searcher: idspace.TransformIDs(imageComponentEdgeSearcher, idspace.NewEdgeToParentTransformer()),
			Options:  imageComponentEdgeMappings.OptionsMap,
		},
		{
			Searcher: imageSearcher,
			Options:  mappings.OptionsMap,
		},
	}...)
}
