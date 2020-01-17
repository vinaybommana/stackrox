package resolvers

import (
	"context"
	"strings"
	"time"

	protoTypes "github.com/gogo/protobuf/types"
	"github.com/graph-gophers/graphql-go"
	"github.com/stackrox/rox/central/graphql/resolvers/loaders"
	"github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	pkgMetrics "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/search"
)

// Cve resolves a single vulnerability based on an id (the CVE value).
func (resolver *Resolver) Cve(ctx context.Context, args idQuery) (*cVEResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "CVE")
	if err := readCVEs(ctx); err != nil {
		return nil, err
	}
	return resolver.wrapCVE(resolver.CVEDataStore.Get(ctx, string(*args.ID)))
}

// Cves resolves a set of vulnerabilities based on a query.
func (resolver *Resolver) Cves(ctx context.Context, args PaginatedQuery) ([]*cVEResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "CVEs")
	if err := readCVEs(ctx); err != nil {
		return nil, err
	}

	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	return resolver.wrapCVEs(resolver.CVEDataStore.SearchRawCVEs(ctx, query))
}

// CveCount returns count of all vulnerabilities across infrastructure
func (resolver *Resolver) CveCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "CVECount")
	if err := readCVEs(ctx); err != nil {
		return 0, err
	}

	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}

	results, err := resolver.CVEDataStore.Search(ctx, q)
	if err != nil {
		return 0, err
	}
	return int32(len(results)), nil
}

// TODO: this and below should only be accessible from a single component context.
// FixedByVersion returns the version of the parent component that removes this CVE.
func (cver *cVEResolver) FixedByVersion(ctx context.Context) (string, error) {
	if err := readCVEs(ctx); err != nil {
		return "", err
	}

	fixableCVEQuery := search.NewConjunctionQuery(
		search.NewQueryBuilder().AddExactMatches(search.CVE, cver.data.GetId()).ProtoQuery(),
		search.NewQueryBuilder().AddBools(search.Fixable, true).ProtoQuery(),
	)

	edges, err := cver.root.ComponentCVEEdgeDataStore.SearchRawEdges(ctx, fixableCVEQuery)
	if err != nil {
		return "", err
	}

	versions := make([]string, 0, len(edges))
	for _, edge := range edges {
		if edge.HasFixedBy == nil {
			continue
		}
		versions = append(versions, edge.GetFixedBy())
	}

	return strings.Join(versions, ","), nil
}

// IsFixable returns whether vulnerability is fixable by any component.
func (cver *cVEResolver) IsFixable(ctx context.Context) (bool, error) {
	if err := readCVEs(ctx); err != nil {
		return false, err
	}

	fixableCVEQuery := search.NewConjunctionQuery(
		search.NewQueryBuilder().AddExactMatches(search.CVE, cver.data.GetId()).ProtoQuery(),
		search.NewQueryBuilder().AddBools(search.Fixable, true).ProtoQuery(),
	)

	results, err := cver.root.ComponentCVEEdgeDataStore.Search(ctx, fixableCVEQuery)
	if err != nil {
		return false, err
	}

	return len(results) != 0, nil
}

// LastScanned is the last time the vulnerability was scanned in an image.
func (cver *cVEResolver) LastScanned(ctx context.Context) (*graphql.Time, error) {
	lastScanned, err := cver.getLastScannedTime(ctx)
	if err != nil {
		return nil, err
	}
	return timestamp(lastScanned)
}

func (cver *cVEResolver) getLastScannedTime(ctx context.Context) (*protoTypes.Timestamp, error) {
	imageResolver, err := cver.loadImages(ctx, search.EmptyQuery())
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

func (cver *cVEResolver) loadImages(ctx context.Context, query *v1.Query) ([]*imageResolver, error) {
	imageLoader, err := loaders.GetImageLoader(ctx)
	if err != nil {
		return nil, err
	}

	pagination := query.GetPagination()
	query.Pagination = nil

	query, err = search.AddAsConjunction(cver.cveQuery(), query)
	if err != nil {
		return nil, err
	}

	query.Pagination = pagination

	return cver.root.wrapImages(imageLoader.FromQuery(ctx, query))
}

func (cver *cVEResolver) cveQuery() *v1.Query {
	return search.NewQueryBuilder().AddExactMatches(search.CVE, cver.data.GetId()).ProtoQuery()
}
