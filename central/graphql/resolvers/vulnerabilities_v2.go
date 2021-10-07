package resolvers

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/cve/converter"
	distroctx "github.com/stackrox/rox/central/graphql/resolvers/distroctx"
	"github.com/stackrox/rox/central/graphql/resolvers/loaders"
	"github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/cve"
	"github.com/stackrox/rox/pkg/dackbox/edges"
	"github.com/stackrox/rox/pkg/features"
	pkgMetrics "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/predicate"
	"github.com/stackrox/rox/pkg/search/scoped"
	"github.com/stackrox/rox/pkg/stringutils"
)

// V2 Connections to root.
//////////////////////////

var (
	cvePredicateFactory        = predicate.NewFactory("cve", &storage.CVE{})
	cvePostFilteringOptionsMap = func() search.OptionsMap {
		opts := search.Walk(v1.SearchCategory_VULNERABILITIES, "cve", (*storage.CVE)(nil))

		cvss := opts.MustGet(search.CVSS.String())
		severity := opts.MustGet(search.Severity.String())

		return search.NewOptionsMap(v1.SearchCategory_VULNERABILITIES).
			Add(search.CVSS, cvss).
			Add(search.Severity, severity)
	}()
)

func getImageIDFromQuery(q *v1.Query) string {
	if q == nil {
		return ""
	}
	var imageID string
	search.ApplyFnToAllBaseQueries(q, func(bq *v1.BaseQuery) {
		matchFieldQuery, ok := bq.GetQuery().(*v1.BaseQuery_MatchFieldQuery)
		if !ok {
			return
		}
		if strings.EqualFold(matchFieldQuery.MatchFieldQuery.GetField(), search.ImageSHA.String()) {
			imageID = matchFieldQuery.MatchFieldQuery.Value
			imageID = strings.TrimRight(imageID, `"`)
			imageID = strings.TrimLeft(imageID, `"`)
		}
	})
	return imageID
}

// AddDistroContext adds the image distribution from the query or scope query if necessary
func (resolver *Resolver) AddDistroContext(ctx context.Context, query, scopeQuery *v1.Query) (context.Context, error) {
	if distro := distroctx.FromContext(ctx); distro != "" {
		return ctx, nil
	}

	scope, ok := scoped.GetScope(ctx)
	if ok && scope.Level == v1.SearchCategory_IMAGES {
		if image := resolver.getImage(ctx, scope.ID); image != nil {
			return distroctx.Context(ctx, image.GetScan().GetOperatingSystem()), nil
		}
	}

	imageIDFromQuery := getImageIDFromQuery(query)
	imageIDFromScope := getImageIDFromQuery(scopeQuery)

	if imageID := stringutils.FirstNonEmpty(imageIDFromQuery, imageIDFromScope); imageID != "" {
		image, exists, err := resolver.ImageDataStore.GetImageMetadata(ctx, imageID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return ctx, nil
		}
		return distroctx.Context(ctx, image.GetScan().GetOperatingSystem()), nil
	}
	return ctx, nil
}

func filterNamespacedFields(query *v1.Query, cves []*storage.CVE) ([]*storage.CVE, error) {
	vulnQuery, _ := search.FilterQueryWithMap(query, cvePostFilteringOptionsMap)
	vulnPred, err := cvePredicateFactory.GeneratePredicate(vulnQuery)
	if err != nil {
		return nil, err
	}
	filtered := cves[:0]
	for _, cve := range cves {
		if vulnPred.Matches(cve) {
			filtered = append(filtered, cve)
		}
	}
	return filtered, nil
}

func needsPostSorting(query *v1.Query) bool {
	for _, so := range query.GetPagination().GetSortOptions() {
		switch so.GetField() {
		case search.Severity.String(), search.CVSS.String():
			return true
		default:
			return false
		}
	}
	return false
}

func sortNamespacedFields(query *v1.Query, cves []*storage.CVE) ([]*storage.CVE, error) {
	// Currently, only one sort option is supported on this endpoint
	sortOption := query.GetPagination().SortOptions[0]
	switch sortOption.Field {
	case search.Severity.String():
		sort.Slice(cves, func(i, j int) bool {
			var result bool
			if cves[i].GetSeverity() != cves[j].GetSeverity() {
				result = cves[i].GetSeverity() < cves[j].GetSeverity()
			} else {
				result = cves[i].GetCvss() < cves[j].GetCvss()
			}
			if sortOption.GetReversed() {
				return !result
			}
			return result
		})
	case search.CVSS.String():
		sort.Slice(cves, func(i, j int) bool {
			var result bool
			if cves[i].GetCvss() != cves[j].GetCvss() {
				result = cves[i].GetCvss() < cves[j].GetCvss()
			} else {
				result = cves[i].GetSeverity() < cves[j].GetSeverity()
			}
			if sortOption.Reversed {
				return !result
			}
			return result
		})
	}
	return cves, nil
}

func (resolver *cVEResolver) Cvss(ctx context.Context) float64 {
	value := resolver.data.GetCvss()
	return float64(value)
}

func (resolver *cVEResolver) CvssV2(ctx context.Context) (*cVSSV2Resolver, error) {
	value := resolver.data.GetCvssV2()
	return resolver.root.wrapCVSSV2(value, true, nil)
}

func (resolver *cVEResolver) CvssV3(ctx context.Context) (*cVSSV3Resolver, error) {
	value := resolver.data.GetCvssV3()
	return resolver.root.wrapCVSSV3(value, true, nil)
}

func (resolver *Resolver) vulnerabilityV2(ctx context.Context, args IDQuery) (VulnerabilityResolver, error) {
	if err := readCVEs(ctx); err != nil {
		return nil, err
	}
	vuln, exists, err := resolver.CVEDataStore.Get(ctx, string(*args.ID))
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, errors.Errorf("cve not found: %s", string(*args.ID))
	}
	vulnResolver, err := resolver.wrapCVE(vuln, true, nil)
	if err != nil {
		return nil, err
	}
	vulnResolver.ctx = ctx
	return vulnResolver, nil
}

func (resolver *Resolver) vulnerabilitiesV2(ctx context.Context, args PaginatedQuery) ([]VulnerabilityResolver, error) {
	if err := readCVEs(ctx); err != nil {
		return nil, err
	}
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	scopeQuery, err := args.AsV1ScopeQueryOrEmpty()
	if err != nil {
		return nil, err
	}

	ctx, err = resolver.AddDistroContext(ctx, query, scopeQuery)
	if err != nil {
		return nil, err
	}
	return resolver.vulnerabilitiesV2Query(ctx, query)
}

func (resolver *Resolver) vulnerabilitiesV2Query(ctx context.Context, query *v1.Query) ([]VulnerabilityResolver, error) {
	vulnLoader, err := loaders.GetCVELoader(ctx)
	if err != nil {
		return nil, err
	}

	query = tryUnsuppressedQuery(query)

	originalQuery := query.Clone()
	var queryModified, postSortingNeeded bool
	if distroctx.IsImageScoped(ctx) {
		query, queryModified = search.InverseFilterQueryWithMap(query, cvePostFilteringOptionsMap) // CVE queryModified
		postSortingNeeded = needsPostSorting(originalQuery)
		// We remove pagination since we want to ensure that result is correct by pushing the pagination to happen after the post sorting.
		if postSortingNeeded {
			query.Pagination = nil
		}
	}

	vulns, err := vulnLoader.FromQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	if queryModified {
		vulns, err = filterNamespacedFields(originalQuery, vulns)
		if err != nil {
			return nil, err
		}
	}
	if postSortingNeeded {
		vulns, err = sortNamespacedFields(originalQuery, vulns)
		if err != nil {
			return nil, err
		}
	}

	// If query was modified, it means the result was not paginated since the filtering removes pagination.
	// If post sorting was needed, which means pagination was not performed because it was it was removed above.
	if queryModified || postSortingNeeded {
		paginatedVulns, err := paginationWrapper{
			pv: originalQuery.GetPagination(),
		}.paginate(vulns, nil)
		if err != nil {
			return nil, err
		}
		vulns = paginatedVulns.([]*storage.CVE)
	}

	vulnResolvers, err := resolver.wrapCVEs(vulns, err)
	ret := make([]VulnerabilityResolver, 0, len(vulns))
	for _, resolver := range vulnResolvers {
		resolver.ctx = ctx
		ret = append(ret, resolver)
	}
	return ret, err
}

func (resolver *Resolver) vulnerabilityCountV2(ctx context.Context, args RawQuery) (int32, error) {
	if err := readCVEs(ctx); err != nil {
		return 0, err
	}
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}

	return resolver.vulnerabilityCountV2Query(ctx, query)
}

func (resolver *Resolver) vulnerabilityCountV2Query(ctx context.Context, query *v1.Query) (int32, error) {
	vulnLoader, err := loaders.GetCVELoader(ctx)
	if err != nil {
		return 0, err
	}

	if distroctx.IsImageScoped(ctx) {
		_, queryModified := search.InverseFilterQueryWithMap(query, cvePostFilteringOptionsMap)
		if queryModified {
			vulns, err := resolver.vulnerabilitiesV2Query(ctx, query)
			if err != nil {
				return 0, err
			}
			return int32(len(vulns)), nil
		}
	}

	query = tryUnsuppressedQuery(query)
	return vulnLoader.CountFromQuery(ctx, query)
}

func (resolver *Resolver) vulnCounterV2(ctx context.Context, args RawQuery) (*VulnerabilityCounterResolver, error) {
	if err := readCVEs(ctx); err != nil {
		return nil, err
	}
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}
	return resolver.vulnCounterV2Query(ctx, query)
}

func (resolver *Resolver) vulnCounterV2Query(ctx context.Context, query *v1.Query) (*VulnerabilityCounterResolver, error) {
	vulnLoader, err := loaders.GetCVELoader(ctx)
	if err != nil {
		return nil, err
	}
	query = tryUnsuppressedQuery(query)
	fixableVulnsQuery := search.ConjunctionQuery(query, search.NewQueryBuilder().AddBools(search.Fixable, true).ProtoQuery())
	fixableVulns, err := vulnLoader.FromQuery(ctx, fixableVulnsQuery)
	if err != nil {
		return nil, err
	}

	unFixableVulnsQuery := search.ConjunctionQuery(query, search.NewQueryBuilder().AddBools(search.Fixable, false).ProtoQuery())
	unFixableCVEs, err := vulnLoader.FromQuery(ctx, unFixableVulnsQuery)
	if err != nil {
		return nil, err
	}

	return mapCVEsToVulnerabilityCounter(fixableVulns, unFixableCVEs), nil
}

func (resolver *Resolver) k8sVulnerabilityV2(ctx context.Context, args IDQuery) (VulnerabilityResolver, error) {
	return resolver.vulnerabilityV2(ctx, args)
}

func (resolver *Resolver) k8sVulnerabilitiesV2(ctx context.Context, q PaginatedQuery) ([]VulnerabilityResolver, error) {
	query := search.AddRawQueriesAsConjunction(q.String(),
		search.NewQueryBuilder().AddExactMatches(search.CVEType, storage.CVE_K8S_CVE.String()).Query())
	return resolver.vulnerabilitiesV2(ctx, PaginatedQuery{Query: &query, Pagination: q.Pagination})
}

func (resolver *Resolver) istioVulnerabilityV2(ctx context.Context, args IDQuery) (VulnerabilityResolver, error) {
	return resolver.vulnerabilityV2(ctx, args)
}

func (resolver *Resolver) istioVulnerabilitiesV2(ctx context.Context, q PaginatedQuery) ([]VulnerabilityResolver, error) {
	query := search.AddRawQueriesAsConjunction(q.String(),
		search.NewQueryBuilder().AddExactMatches(search.CVEType, storage.CVE_ISTIO_CVE.String()).Query())
	return resolver.vulnerabilitiesV2(ctx, PaginatedQuery{Query: &query, Pagination: q.Pagination})
}

func (resolver *Resolver) openShiftVulnerabilityV2(ctx context.Context, args IDQuery) (VulnerabilityResolver, error) {
	return resolver.vulnerabilityV2(ctx, args)
}

func (resolver *Resolver) openShiftVulnerabilitiesV2(ctx context.Context, q PaginatedQuery) ([]VulnerabilityResolver, error) {
	query := search.AddRawQueriesAsConjunction(q.String(),
		search.NewQueryBuilder().AddExactMatches(search.CVEType, storage.CVE_OPENSHIFT_CVE.String()).Query())
	return resolver.vulnerabilitiesV2(ctx, PaginatedQuery{Query: &query, Pagination: q.Pagination})
}

// Implemented Resolver.
////////////////////////

func (resolver *cVEResolver) ID(ctx context.Context) graphql.ID {
	value := resolver.data.GetId()
	return graphql.ID(value)
}

func (resolver *cVEResolver) Cve(ctx context.Context) string {
	return resolver.data.GetId()
}

func (resolver *cVEResolver) getCVEQuery() *v1.Query {
	return search.NewQueryBuilder().AddExactMatches(search.CVE, resolver.data.GetId()).ProtoQuery()
}

// IsFixable returns whether vulnerability is fixable by any component.
func (resolver *cVEResolver) IsFixable(_ context.Context, args RawQuery) (bool, error) {
	// CVE is used in scoping but it's not relevant to IsFixable because it is already scoped to a CVE
	q, err := args.AsV1QueryOrEmpty(search.ExcludeFieldLabel(search.CVE))
	if err != nil {
		return false, err
	}

	conjuncts := []*v1.Query{q}

	ctx := resolver.ctx
	if scope, ok := scoped.GetScope(ctx); !ok {
		ctx = resolver.scopeContext(ctx)
	} else if scope.Level != v1.SearchCategory_VULNERABILITIES {
		// If the scope is not set to vulnerabilities then
		// we need to add a query to scope the search to the current vuln
		conjuncts = append(conjuncts, resolver.getCVEQuery())
	}

	if cve.ContainsComponentBasedCVE(resolver.data.GetTypes()) {
		query := search.ConjunctionQuery(append(conjuncts, search.NewQueryBuilder().AddBools(search.Fixable, true).ProtoQuery())...)
		count, err := resolver.root.ComponentCVEEdgeDataStore.Count(ctx, query)
		if err != nil {
			return false, err
		}
		if count != 0 {
			return true, nil
		}
	}
	if cve.ContainsClusterCVE(resolver.data.GetTypes()) {
		query := search.ConjunctionQuery(append(conjuncts, search.NewQueryBuilder().AddBools(search.ClusterCVEFixable, true).ProtoQuery())...)
		count, err := resolver.root.clusterCVEEdgeDataStore.Count(ctx, query)
		if err != nil {
			return false, err
		}
		if count != 0 {
			return true, nil
		}
	}
	return false, nil
}

func (resolver *cVEResolver) scopeContext(ctx context.Context) context.Context {
	return scoped.Context(ctx, scoped.Scope{
		ID:    resolver.data.GetId(),
		Level: v1.SearchCategory_VULNERABILITIES,
	})
}

func (resolver *cVEResolver) getEnvImpactComponentsForImages(ctx context.Context) (numerator, denominator int, err error) {
	allDepsCount, err := resolver.root.DeploymentDataStore.CountDeployments(ctx)
	if err != nil {
		return 0, 0, err
	}
	if allDepsCount == 0 {
		return 0, 0, nil
	}
	deploymentLoader, err := loaders.GetDeploymentLoader(ctx)
	if err != nil {
		return 0, 0, err
	}
	withThisCVECount, err := deploymentLoader.CountFromQuery(resolver.scopeContext(ctx), search.EmptyQuery())
	if err != nil {
		return 0, 0, err
	}
	return int(withThisCVECount), allDepsCount, nil
}

func (resolver *cVEResolver) getEnvImpactComponentsForNodes(ctx context.Context) (numerator, denominator int, err error) {
	allNodesCount, err := resolver.root.NodeGlobalDataStore.CountAllNodes(ctx)
	if err != nil {
		return 0, 0, err
	}
	if allNodesCount == 0 {
		return 0, 0, nil
	}
	nodeLoader, err := loaders.GetNodeLoader(ctx)
	if err != nil {
		return 0, 0, err
	}
	withThisCVECount, err := nodeLoader.CountFromQuery(resolver.scopeContext(ctx), search.EmptyQuery())
	if err != nil {
		return 0, 0, err
	}
	return int(withThisCVECount), allNodesCount, nil
}

// EnvImpact is the fraction of deployments that contains the CVE
func (resolver *cVEResolver) EnvImpact(ctx context.Context) (float64, error) {
	var numerator, denominator int

	for _, vulnType := range resolver.data.GetTypes() {
		var n, d int
		var err error

		switch vulnType {
		case storage.CVE_K8S_CVE:
			n, d, err = resolver.getEnvImpactComponentsForPerClusterVuln(ctx, converter.K8s)
		case storage.CVE_ISTIO_CVE:
			n, d, err = resolver.getEnvImpactComponentsForPerClusterVuln(ctx, converter.Istio)
		case storage.CVE_OPENSHIFT_CVE:
			n, d, err = resolver.getEnvImpactComponentsForPerClusterVuln(ctx, converter.OpenShift)
		case storage.CVE_IMAGE_CVE:
			n, d, err = resolver.getEnvImpactComponentsForImages(ctx)
		case storage.CVE_NODE_CVE:
			n, d, err = resolver.getEnvImpactComponentsForNodes(ctx)
		default:
			return 0, errors.Errorf("unknown CVE type: %s", vulnType)
		}

		if err != nil {
			return 0, err
		}

		numerator += n
		denominator += d
	}

	if denominator == 0 {
		return 0, nil
	}

	return float64(numerator) / float64(denominator), nil
}

func (resolver *cVEResolver) getEnvImpactComponentsForPerClusterVuln(ctx context.Context, ct converter.CVEType) (int, int, error) {
	clusters, err := resolver.root.ClusterDataStore.GetClusters(ctx)
	if err != nil {
		return 0, 0, err
	}
	affectedClusters, err := resolver.root.orchestratorIstioCVEManager.GetAffectedClusters(ctx, resolver.data.GetId(), ct, resolver.root.cveMatcher)
	if err != nil {
		return 0, 0, err
	}
	return len(affectedClusters), len(clusters), nil
}

// LastScanned is the last time the vulnerability was scanned in an image.
func (resolver *cVEResolver) LastScanned(ctx context.Context) (*graphql.Time, error) {
	imageLoader, err := loaders.GetImageLoader(ctx)
	if err != nil {
		return nil, err
	}

	q := search.EmptyQuery()
	q.Pagination = &v1.QueryPagination{
		Limit:  1,
		Offset: 0,
		SortOptions: []*v1.QuerySortOption{
			{
				Field:    search.ImageScanTime.String(),
				Reversed: true,
			},
		},
	}

	images, err := imageLoader.FromQuery(resolver.scopeContext(ctx), q)
	if err != nil || len(images) == 0 {
		return nil, err
	} else if len(images) > 1 {
		return nil, errors.New("multiple images matched for last scanned vulnerability query")
	}

	return timestamp(images[0].GetScan().GetScanTime())
}

func (resolver *cVEResolver) Vectors() *EmbeddedVulnerabilityVectorsResolver {
	if val := resolver.data.GetCvssV3(); val != nil {
		return &EmbeddedVulnerabilityVectorsResolver{
			resolver: &cVSSV3Resolver{resolver.ctx, resolver.root, val},
		}
	}
	if val := resolver.data.GetCvssV2(); val != nil {
		return &EmbeddedVulnerabilityVectorsResolver{
			resolver: &cVSSV2Resolver{resolver.ctx, resolver.root, val},
		}
	}
	return nil
}

func (resolver *cVEResolver) VulnerabilityType() string {
	return resolver.data.GetType().String()
}

func (resolver *cVEResolver) VulnerabilityTypes() []string {
	vulnTypes := make([]string, 0, len(resolver.data.GetTypes()))
	for _, vulnType := range resolver.data.GetTypes() {
		vulnTypes = append(vulnTypes, vulnType.String())
	}
	return vulnTypes
}

// Components are the components that contain the CVE/Vulnerability.
func (resolver *cVEResolver) Components(ctx context.Context, args PaginatedQuery) ([]ComponentResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.CVEs, "Components")

	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	cveQuery := search.NewQueryBuilder().AddExactMatches(search.CVE, resolver.data.GetId()).ProtoQuery()
	query, err = search.AddAsConjunction(cveQuery, query)
	if err != nil {
		return nil, err
	}

	return resolver.root.componentsV2Query(resolver.scopeContext(ctx), query)
}

// ComponentCount is the number of components that contain the CVE/Vulnerability.
func (resolver *cVEResolver) ComponentCount(ctx context.Context, args RawQuery) (int32, error) {
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}
	cveQuery := search.NewQueryBuilder().AddExactMatches(search.CVE, resolver.data.GetId()).ProtoQuery()
	query, err = search.AddAsConjunction(cveQuery, query)
	if err != nil {
		return 0, err
	}

	componentLoader, err := loaders.GetComponentLoader(ctx)
	if err != nil {
		return 0, err
	}
	return componentLoader.CountFromQuery(resolver.scopeContext(ctx), query)
}

// Images are the images that contain the CVE/Vulnerability.
func (resolver *cVEResolver) Images(ctx context.Context, args PaginatedQuery) ([]*imageResolver, error) {
	if err := readImages(ctx); err != nil {
		return []*imageResolver{}, nil
	}

	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	imageLoader, err := loaders.GetImageLoader(ctx)
	if err != nil {
		return nil, err
	}
	return resolver.root.wrapImages(imageLoader.FromQuery(resolver.scopeContext(ctx), query))
}

// ImageCount is the number of images that contain the CVE/Vulnerability.
func (resolver *cVEResolver) ImageCount(ctx context.Context, args RawQuery) (int32, error) {
	if err := readImages(ctx); err != nil {
		return 0, nil
	}
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}
	imageLoader, err := loaders.GetImageLoader(ctx)
	if err != nil {
		return 0, err
	}
	return imageLoader.CountFromQuery(resolver.scopeContext(ctx), query)
}

// Deployments are the deployments that contain the CVE/Vulnerability.
func (resolver *cVEResolver) Deployments(ctx context.Context, args PaginatedQuery) ([]*deploymentResolver, error) {
	if err := readDeployments(ctx); err != nil {
		return nil, err
	}
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	deploymentLoader, err := loaders.GetDeploymentLoader(ctx)
	if err != nil {
		return nil, err
	}
	return resolver.root.wrapDeployments(deploymentLoader.FromQuery(resolver.scopeContext(ctx), query))
}

// DeploymentCount is the number of deployments that contain the CVE/Vulnerability.
func (resolver *cVEResolver) DeploymentCount(ctx context.Context, args RawQuery) (int32, error) {
	if err := readDeployments(ctx); err != nil {
		return 0, err
	}
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}
	deploymentLoader, err := loaders.GetDeploymentLoader(ctx)
	if err != nil {
		return 0, err
	}

	return deploymentLoader.CountFromQuery(resolver.scopeContext(ctx), query)
}

// Nodes are the nodes that contain the CVE/Vulnerability.
func (resolver *cVEResolver) Nodes(ctx context.Context, args PaginatedQuery) ([]*nodeResolver, error) {
	if err := readNodes(ctx); err != nil {
		return []*nodeResolver{}, nil
	}
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	nodeLoader, err := loaders.GetNodeLoader(ctx)
	if err != nil {
		return nil, err
	}
	return resolver.root.wrapNodes(nodeLoader.FromQuery(resolver.scopeContext(ctx), query))
}

// NodeCount is the number of nodes that contain the CVE/Vulnerability.
func (resolver *cVEResolver) NodeCount(ctx context.Context, args RawQuery) (int32, error) {
	if err := readNodes(ctx); err != nil {
		return 0, nil
	}
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}
	nodeLoader, err := loaders.GetNodeLoader(ctx)
	if err != nil {
		return 0, err
	}

	return nodeLoader.CountFromQuery(resolver.scopeContext(ctx), query)
}

// These return dummy values, as they should not be accessed from the top level vuln resolver, but the embedded
// version instead.

// FixedByVersion returns the version of the parent component that removes this CVE.
func (resolver *cVEResolver) FixedByVersion(ctx context.Context) (string, error) {
	return resolver.getCVEFixedByVersion(ctx)
}

// UnusedVarSink represents a query sink
func (resolver *cVEResolver) UnusedVarSink(ctx context.Context, args RawQuery) *int32 {
	return nil
}

func (resolver *cVEResolver) getCVEFixedByVersion(ctx context.Context) (string, error) {
	if cve.ContainsComponentBasedCVE(resolver.data.GetTypes()) {
		return resolver.getComponentFixedByVersion(ctx)
	}
	return resolver.getClusterFixedByVersion(ctx)
}

func (resolver *cVEResolver) getComponentFixedByVersion(_ context.Context) (string, error) {
	scope, hasScope := scoped.GetScope(resolver.ctx)
	if !hasScope {
		return "", nil
	}
	if scope.Level != v1.SearchCategory_IMAGE_COMPONENTS {
		return "", nil
	}

	edgeID := edges.EdgeID{ParentID: scope.ID, ChildID: resolver.data.GetId()}.ToString()
	edge, found, err := resolver.root.ComponentCVEEdgeDataStore.Get(resolver.ctx, edgeID)
	if err != nil || !found {
		return "", err
	}
	return edge.GetFixedBy(), nil
}

func (resolver *cVEResolver) getClusterFixedByVersion(_ context.Context) (string, error) {
	scope, hasScope := scoped.GetScope(resolver.ctx)
	if !hasScope {
		return "", nil
	}
	if scope.Level != v1.SearchCategory_CLUSTERS {
		return "", nil
	}

	edgeID := edges.EdgeID{ParentID: scope.ID, ChildID: resolver.data.GetId()}.ToString()
	edge, found, err := resolver.root.clusterCVEEdgeDataStore.Get(resolver.ctx, edgeID)
	if err != nil || !found {
		return "", err
	}
	return edge.GetFixedBy(), nil
}

func (resolver *cVEResolver) DiscoveredAtImage(ctx context.Context, args RawQuery) (*graphql.Time, error) {
	if !cve.ContainsCVEType(resolver.data.GetTypes(), storage.CVE_IMAGE_CVE) {
		return nil, nil
	}

	var imageID string
	scope, hasScope := scoped.GetScope(resolver.ctx)
	if hasScope && scope.Level == v1.SearchCategory_IMAGES {
		imageID = scope.ID
	} else if !hasScope || scope.Level != v1.SearchCategory_IMAGES {
		var err error
		imageID, err = getImageIDFromIfImageShaQuery(ctx, resolver.root, args)
		if err != nil {
			return nil, errors.Wrap(err, "could not determine vulnerability discovered time in image")
		}
	}

	if imageID == "" {
		return nil, nil
	}

	edgeID := edges.EdgeID{
		ParentID: imageID,
		ChildID:  resolver.data.GetId(),
	}.ToString()

	edge, found, err := resolver.root.ImageCVEEdgeDataStore.Get(resolver.ctx, edgeID)
	if err != nil || !found {
		return nil, err
	}
	return timestamp(edge.GetFirstImageOccurrence())
}

// ActiveState shows the activeness of a vulnerability in a deployment context.
func (resolver *cVEResolver) ActiveState(ctx context.Context, args RawQuery) (*activeStateResolver, error) {
	if !features.ActiveVulnManagement.Enabled() {
		return nil, nil
	}
	scopeQuery, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}
	deploymentID := getDeploymentScope(scopeQuery, ctx, resolver.ctx)
	if deploymentID == "" {
		return nil, nil
	}
	// We only support OS level component. The active state is not determined if there is no OS level component associate with this vuln.
	query := search.NewQueryBuilder().AddExactMatches(search.CVE, resolver.data.GetId()).AddStrings(search.ComponentSource, storage.SourceType_OS.String()).ProtoQuery()
	osLevelComponents, err := resolver.root.ImageComponentDataStore.Count(ctx, query)
	if err != nil {
		return nil, err
	}
	if osLevelComponents == 0 {
		return &activeStateResolver{root: resolver.root, state: Undetermined}, nil
	}

	query = search.ConjunctionQuery(resolver.getCVEQuery(), search.NewQueryBuilder().AddExactMatches(search.DeploymentID, deploymentID).ProtoQuery())
	results, err := resolver.root.ActiveComponent.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	ids := search.ResultsToIDs(results)
	state := Inactive
	if len(ids) != 0 {
		state = Active
	}
	return &activeStateResolver{root: resolver.root, state: state, activeComponentIDs: ids}, nil
}
