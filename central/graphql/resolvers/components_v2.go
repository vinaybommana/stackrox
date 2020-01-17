package resolvers

import (
	"context"
	"time"

	protoTypes "github.com/gogo/protobuf/types"
	"github.com/graph-gophers/graphql-go"
	"github.com/stackrox/rox/central/graphql/resolvers/loaders"
	"github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	pkgMetrics "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/search"
)

// ComponentV2 resolves a single vulnerability based on an id (the CVE value).
func (resolver *Resolver) ComponentV2(ctx context.Context, args idQuery) (*imageComponentResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "Component")
	if err := readImageComponents(ctx); err != nil {
		return nil, err
	}
	return resolver.wrapImageComponent(resolver.ImageComponentDataStore.Get(ctx, string(*args.ID)))
}

// ComponentsV2 resolves a set of vulnerabilities based on a query.
func (resolver *Resolver) ComponentsV2(ctx context.Context, args PaginatedQuery) ([]*imageComponentResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "Components")
	if err := readImageComponents(ctx); err != nil {
		return nil, err
	}

	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}
	return resolver.wrapImageComponents(resolver.ImageComponentDataStore.SearchRawImageComponents(ctx, query))
}

// ComponentCountV2 returns count of all vulnerabilities across infrastructure
func (resolver *Resolver) ComponentCountV2(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "ComponentCount")
	if err := readImageComponents(ctx); err != nil {
		return 0, err
	}

	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}

	results, err := resolver.ImageComponentDataStore.Search(ctx, q)
	if err != nil {
		return 0, err
	}
	return int32(len(results)), nil
}

// LastScanned is the last time the component was scanned in an image.
func (icr *imageComponentResolver) LastScanned(ctx context.Context) (*graphql.Time, error) {
	lastScanned, err := icr.getLastScannedTime(ctx)
	if err != nil {
		return nil, err
	}

	return timestamp(lastScanned)
}

func (icr *imageComponentResolver) getLastScannedTime(ctx context.Context) (*protoTypes.Timestamp, error) {
	imageResolver, err := icr.loadImages(ctx, search.EmptyQuery())
	if err != nil {
		return nil, err
	}

	var lastScanned *protoTypes.Timestamp
	for _, resolver := range imageResolver {
		if resolver.data.GetScan() == nil {
			continue
		}

		if lastScanned == nil || resolver.data.GetScan().GetScanTime().Compare(lastScanned) > 0 {
			lastScanned = resolver.data.GetScan().GetScanTime()
		}
	}

	return lastScanned, nil
}

func (icr *imageComponentResolver) loadImages(ctx context.Context, query *v1.Query) ([]*imageResolver, error) {
	imageLoader, err := loaders.GetImageLoader(ctx)
	if err != nil {
		return nil, err
	}

	pagination := query.GetPagination()
	query.Pagination = nil

	query, err = search.AddAsConjunction(icr.componentQuery(), query)
	if err != nil {
		return nil, err
	}

	query.Pagination = pagination

	return icr.root.wrapImages(imageLoader.FromQuery(ctx, query))
}

func (icr *imageComponentResolver) componentQuery() *v1.Query {
	return search.NewQueryBuilder().
		AddExactMatches(search.Component, icr.data.GetName()).
		AddExactMatches(search.ComponentVersion, icr.data.GetVersion()).
		ProtoQuery()
}
