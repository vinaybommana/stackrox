package resolvers

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	protoTypes "github.com/gogo/protobuf/types"
	"github.com/graph-gophers/graphql-go"
	"github.com/stackrox/rox/central/cve/converter"
	"github.com/stackrox/rox/central/graphql/resolvers/loaders"
	"github.com/stackrox/rox/central/image/mappings"
	imageComponentConverter "github.com/stackrox/rox/central/imagecomponent/converter"
	"github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	pkgMetrics "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/predicate"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	componentPredicateFactory = predicate.NewFactory("component", &storage.EmbeddedImageScanComponent{})
)

func init() {
	schema := getBuilder()
	utils.Must(
		schema.AddType("EmbeddedImageScanComponent", []string{
			"license: License",
			"id: ID!",
			"name: String!",
			"version: String!",
			"topVuln: EmbeddedVulnerability",
			"vulns(query: String, pagination: Pagination): [EmbeddedVulnerability]!",
			"vulnCount(query: String): Int!",
			"vulnCounter(query: String): VulnerabilityCounter!",
			"lastScanned: Time",
			"images(query: String, pagination: Pagination): [Image!]!",
			"imageCount(query: String): Int!",
			"deployments(query: String, pagination: Pagination): [Deployment!]!",
			"deploymentCount(query: String): Int!",
			"priority: Int!",
			"source: String!",
			"location: String!",
		}),
		schema.AddExtraResolver("ImageScan", `components(query: String, pagination: Pagination): [EmbeddedImageScanComponent!]!`),
		schema.AddExtraResolver("ImageScan", `componentCount(query: String): Int!`),
		schema.AddQuery("component(id: ID): EmbeddedImageScanComponent"),
		schema.AddQuery("components(query: String, pagination: Pagination): [EmbeddedImageScanComponent!]!"),
		schema.AddQuery("componentCount(query: String): Int!"),
	)
}

// Component returns an image scan component based on an input id (name:version)
func (resolver *Resolver) Component(ctx context.Context, args idQuery) (*EmbeddedImageScanComponentResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "ImageComponent")
	if features.Dackbox.Enabled() {
		return resolver.getComponentFromDackBox(ctx, args)
	}

	if err := readImages(ctx); err != nil {
		return nil, err
	}

	cID, err := componentIDFromString(string(*args.ID))
	if err != nil {
		return nil, err
	}

	query := search.NewQueryBuilder().
		AddExactMatches(search.Component, cID.Name).
		AddExactMatches(search.ComponentVersion, cID.Version).
		ProtoQuery()
	comps, err := components(ctx, resolver, query)
	if err != nil {
		return nil, err
	} else if len(comps) == 0 {
		return nil, nil
	} else if len(comps) > 1 {
		return nil, fmt.Errorf("multiple components matched: %s this should not happen", string(*args.ID))
	}
	return comps[0], nil
}

// Components returns the image scan components that match the input query.
func (resolver *Resolver) Components(ctx context.Context, q PaginatedQuery) ([]*EmbeddedImageScanComponentResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "ImageComponents")
	if features.Dackbox.Enabled() {
		return resolver.getComponentsFromDackBox(ctx, q)
	}

	if err := readImages(ctx); err != nil {
		return nil, err
	}

	// Convert to query, but link the fields for the search.
	query, err := q.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	resolvers, err := paginationWrapper{
		pv: query.Pagination,
	}.paginate(components(ctx, resolver, query))
	return resolvers.([]*EmbeddedImageScanComponentResolver), err
}

// ComponentCount returns count of all clusters across infrastructure
func (resolver *Resolver) ComponentCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "ComponentCount")
	if features.Dackbox.Enabled() {
		return resolver.ComponentCountV2(ctx, args)
	}

	if err := readImages(ctx); err != nil {
		return 0, err
	}
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}
	comps, err := components(ctx, resolver, query)
	if err != nil {
		return 0, err
	}
	return int32(len(comps)), nil
}

// Helper function that actually runs the queries and produces the resolvers from the images.
func components(ctx context.Context, root *Resolver, query *v1.Query) ([]*EmbeddedImageScanComponentResolver, error) {
	imageLoader, err := loaders.GetImageLoader(ctx)
	if err != nil {
		return nil, err
	}

	// Run search on images.
	images, err := imageLoader.FromQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	return mapImagesToComponentResolvers(root, images, query)
}

func (resolver *imageScanResolver) Components(ctx context.Context, args PaginatedQuery) ([]*EmbeddedImageScanComponentResolver, error) {
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	pagination := query.GetPagination()
	query.Pagination = nil

	vulns, err := mapImagesToComponentResolvers(resolver.root, []*storage.Image{
		{
			Scan: resolver.data,
		},
	}, query)

	resolvers, err := paginationWrapper{
		pv: pagination,
	}.paginate(vulns, err)
	return resolvers.([]*EmbeddedImageScanComponentResolver), err
}

func (resolver *imageScanResolver) ComponentCount(ctx context.Context, args RawQuery) (int32, error) {
	resolvers, err := resolver.Components(ctx, PaginatedQuery{Query: args.Query})
	if err != nil {
		return 0, err
	}
	return int32(len(resolvers)), nil
}

// EmbeddedImageScanComponentResolver resolves data about an image scan component.
type EmbeddedImageScanComponentResolver struct {
	root        *Resolver
	lastScanned *protoTypes.Timestamp
	data        *storage.EmbeddedImageScanComponent
}

// License return the license for the image component.
func (eicr *EmbeddedImageScanComponentResolver) License(ctx context.Context) (*licenseResolver, error) {
	value := eicr.data.GetLicense()
	return eicr.root.wrapLicense(value, true, nil)
}

// ID returns a unique identifier for the component.
func (eicr *EmbeddedImageScanComponentResolver) ID(ctx context.Context) graphql.ID {
	cID := &componentID{
		Name:    eicr.data.GetName(),
		Version: eicr.data.GetVersion(),
	}
	return graphql.ID(cID.toString())
}

// Name returns the name of the component.
func (eicr *EmbeddedImageScanComponentResolver) Name(ctx context.Context) string {
	return eicr.data.GetName()
}

// Version gives the version of the image component.
func (eicr *EmbeddedImageScanComponentResolver) Version(ctx context.Context) string {
	return eicr.data.GetVersion()
}

// Priority returns the priority of the component.
func (eicr *EmbeddedImageScanComponentResolver) Priority(ctx context.Context) int32 {
	return int32(eicr.data.GetPriority())
}

// Source returns the source of the component.
// TODO: replace this placeholder return value with actual data when available
func (eicr *EmbeddedImageScanComponentResolver) Source(ctx context.Context) string {
	return "placeholder source"
}

// Location returns the location of the component.
// TODO: replace this placeholder return value with actual data when available
func (eicr *EmbeddedImageScanComponentResolver) Location(ctx context.Context) string {
	return "placeholder location"
}

// LayerIndex is the index in the parent image.
// TODO: make this only accessable when coming from an image resolver.
func (eicr *EmbeddedImageScanComponentResolver) LayerIndex() *int32 {
	w, ok := eicr.data.GetHasLayerIndex().(*storage.EmbeddedImageScanComponent_LayerIndex)
	if !ok {
		return nil
	}
	v := w.LayerIndex
	return &v
}

// LastScanned is the last time the vulnerability was scanned in an image.
func (eicr *EmbeddedImageScanComponentResolver) LastScanned(ctx context.Context) (*graphql.Time, error) {
	return timestamp(eicr.lastScanned)
}

// TopVuln returns the first vulnerability with the top CVSS score.
func (eicr *EmbeddedImageScanComponentResolver) TopVuln(ctx context.Context) (*EmbeddedVulnerabilityResolver, error) {
	var maxCvss *storage.EmbeddedVulnerability
	for _, vuln := range eicr.data.GetVulns() {
		if maxCvss == nil || vuln.GetCvss() > maxCvss.GetCvss() {
			maxCvss = vuln
		}
	}
	if maxCvss == nil {
		return nil, nil
	}
	return eicr.root.wrapEmbeddedVulnerability(maxCvss, nil)
}

// Vulns resolves the vulnerabilities contained in the image component.
func (eicr *EmbeddedImageScanComponentResolver) Vulns(ctx context.Context, args PaginatedQuery) ([]*EmbeddedVulnerabilityResolver, error) {
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	vulnQuery, _ := search.FilterQueryWithMap(query, mappings.VulnerabilityOptionsMap)
	vulnPred, err := vulnPredicateFactory.GeneratePredicate(vulnQuery)
	if err != nil {
		return nil, err
	}

	// Use the images to map CVEs to the images and components.
	vulns := make([]*EmbeddedVulnerabilityResolver, 0, len(eicr.data.GetVulns()))
	for _, vuln := range eicr.data.GetVulns() {
		if !vulnPred.Matches(vuln) {
			continue
		}
		vulns = append(vulns, &EmbeddedVulnerabilityResolver{
			data:        vuln,
			root:        eicr.root,
			lastScanned: eicr.lastScanned,
		})
	}

	resolvers, err := paginationWrapper{
		pv: query.GetPagination(),
	}.paginate(vulns, nil)
	return resolvers.([]*EmbeddedVulnerabilityResolver), err
}

// VulnCount resolves the number of vulnerabilities contained in the image component.
func (eicr *EmbeddedImageScanComponentResolver) VulnCount(ctx context.Context, args RawQuery) (int32, error) {
	if features.Dackbox.Enabled() {
		return eicr.root.VulnerabilityCount(ctx, args)
	}

	vulns, err := eicr.Vulns(ctx, PaginatedQuery{Query: args.Query})
	if err != nil {
		return 0, err
	}
	return int32(len(vulns)), nil
}

// VulnCounter resolves the number of different types of vulnerabilities contained in an image component.
func (eicr *EmbeddedImageScanComponentResolver) VulnCounter(ctx context.Context, args RawQuery) (*VulnerabilityCounterResolver, error) {
	vulnResolvers, err := eicr.Vulns(ctx, PaginatedQuery{Query: args.Query})
	if err != nil {
		return nil, err
	}

	vulns := make([]*storage.EmbeddedVulnerability, 0, len(vulnResolvers))
	for _, vulnResolver := range vulnResolvers {
		vulns = append(vulns, vulnResolver.data)
	}

	return mapVulnsToVulnerabilityCounter(vulns), nil
}

// Images are the images that contain the Component.
func (eicr *EmbeddedImageScanComponentResolver) Images(ctx context.Context, args PaginatedQuery) ([]*imageResolver, error) {
	// Convert to query, but link the fields for the search.
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}
	images, err := eicr.loadImages(ctx, query)
	if err != nil {
		return nil, err
	}
	return images, nil
}

// ImageCount is the number of images that contain the Component.
func (eicr *EmbeddedImageScanComponentResolver) ImageCount(ctx context.Context, args RawQuery) (int32, error) {
	imageLoader, err := loaders.GetImageLoader(ctx)
	if err != nil {
		return 0, err
	}
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}
	query, err = search.AddAsConjunction(eicr.componentQuery(), query)
	if err != nil {
		return 0, err
	}
	return imageLoader.CountFromQuery(ctx, query)
}

// Deployments are the deployments that contain the Component.
func (eicr *EmbeddedImageScanComponentResolver) Deployments(ctx context.Context, args PaginatedQuery) ([]*deploymentResolver, error) {
	if err := readDeployments(ctx); err != nil {
		return nil, err
	}
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}
	return eicr.loadDeployments(ctx, query)
}

// DeploymentCount is the number of deployments that contain the Component.
func (eicr *EmbeddedImageScanComponentResolver) DeploymentCount(ctx context.Context, args RawQuery) (int32, error) {
	if err := readDeployments(ctx); err != nil {
		return 0, err
	}
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}
	deploymentBaseQuery, err := eicr.getDeploymentBaseQuery(ctx)
	if err != nil || deploymentBaseQuery == nil {
		return 0, err
	}
	deploymentLoader, err := loaders.GetDeploymentLoader(ctx)
	if err != nil {
		return 0, err
	}
	return deploymentLoader.CountFromQuery(ctx, search.ConjunctionQuery(deploymentBaseQuery, query))
}

func (eicr *EmbeddedImageScanComponentResolver) loadImages(ctx context.Context, query *v1.Query) ([]*imageResolver, error) {
	imageLoader, err := loaders.GetImageLoader(ctx)
	if err != nil {
		return nil, err
	}

	pagination := query.GetPagination()
	query.Pagination = nil

	query, err = search.AddAsConjunction(eicr.componentQuery(), query)
	if err != nil {
		return nil, err
	}

	query.Pagination = pagination

	return eicr.root.wrapImages(imageLoader.FromQuery(ctx, query))
}

func (eicr *EmbeddedImageScanComponentResolver) loadDeployments(ctx context.Context, query *v1.Query) ([]*deploymentResolver, error) {
	deploymentBaseQuery, err := eicr.getDeploymentBaseQuery(ctx)
	if err != nil || deploymentBaseQuery == nil {
		return nil, err
	}

	ListDeploymentLoader, err := loaders.GetListDeploymentLoader(ctx)
	if err != nil {
		return nil, err
	}

	pagination := query.GetPagination()
	query.Pagination = nil

	query, err = search.AddAsConjunction(deploymentBaseQuery, query)
	if err != nil {
		return nil, err
	}

	query.Pagination = pagination

	return eicr.root.wrapListDeployments(ListDeploymentLoader.FromQuery(ctx, query))
}

func (eicr *EmbeddedImageScanComponentResolver) getDeploymentBaseQuery(ctx context.Context) (*v1.Query, error) {
	imageQuery := eicr.componentQuery()
	results, err := eicr.root.ImageDataStore.Search(ctx, imageQuery)
	if err != nil || len(results) == 0 {
		return nil, err
	}

	// Create a query that finds all of the deployments that contain at least one of the infected images.
	var qb []*v1.Query
	for _, id := range search.ResultsToIDs(results) {
		qb = append(qb, search.NewQueryBuilder().AddExactMatches(search.ImageSHA, id).ProtoQuery())
	}
	return search.DisjunctionQuery(qb...), nil
}

func (eicr *EmbeddedImageScanComponentResolver) componentQuery() *v1.Query {
	return search.NewQueryBuilder().
		AddExactMatches(search.Component, eicr.data.GetName()).
		AddExactMatches(search.ComponentVersion, eicr.data.GetVersion()).
		ProtoQuery()
}

// Static helpers.
//////////////////

// Synthetic ID for component objects composed of the name and version of the component.
type componentID struct {
	Name    string
	Version string
}

func componentIDFromString(str string) (*componentID, error) {
	nameAndVersionEncoded := strings.Split(str, ":")
	if len(nameAndVersionEncoded) != 2 {
		return nil, fmt.Errorf("invalid id: %s", str)
	}
	name, err := base64.URLEncoding.DecodeString(nameAndVersionEncoded[0])
	if err != nil {
		return nil, err
	}
	version, err := base64.URLEncoding.DecodeString(nameAndVersionEncoded[1])
	if err != nil {
		return nil, err
	}
	return &componentID{Name: string(name), Version: string(version)}, nil
}

func (cID *componentID) toString() string {
	nameEncoded := base64.URLEncoding.EncodeToString([]byte(cID.Name))
	versionEncoded := base64.URLEncoding.EncodeToString([]byte(cID.Version))
	return fmt.Sprintf("%s:%s", nameEncoded, versionEncoded)
}

// Map the images that matched a query to the image components it contains.
func mapImagesToComponentResolvers(root *Resolver, images []*storage.Image, query *v1.Query) ([]*EmbeddedImageScanComponentResolver, error) {
	query, _ = search.FilterQueryWithMap(query, mappings.ComponentOptionsMap)
	componentPred, err := componentPredicateFactory.GeneratePredicate(query)
	if err != nil {
		return nil, err
	}

	// Use the images to map CVEs to the images and components.
	idToComponent := make(map[componentID]*EmbeddedImageScanComponentResolver)
	for _, image := range images {
		for _, component := range image.GetScan().GetComponents() {
			if !componentPred.Matches(component) {
				continue
			}
			thisComponentID := componentID{Name: component.GetName(), Version: component.GetVersion()}
			if _, exists := idToComponent[thisComponentID]; !exists {
				idToComponent[thisComponentID] = &EmbeddedImageScanComponentResolver{
					root: root,
					data: component,
				}
			}
			latestTime := idToComponent[thisComponentID].lastScanned
			if latestTime == nil || image.GetScan().GetScanTime().Compare(latestTime) > 0 {
				idToComponent[thisComponentID].lastScanned = image.GetScan().GetScanTime()
			}
		}
	}

	// Create the resolvers.
	resolvers := make([]*EmbeddedImageScanComponentResolver, 0, len(idToComponent))
	for _, component := range idToComponent {
		resolvers = append(resolvers, component)
	}
	return resolvers, nil
}

func (resolver *Resolver) getComponentFromDackBox(ctx context.Context, args idQuery) (*EmbeddedImageScanComponentResolver, error) {
	component, err := resolver.ComponentV2(ctx, args)
	if err != nil {
		return nil, err
	}

	cves, err := resolver.getComponentCVEsFromDackBox(ctx, component.data.GetId())
	if err != nil {
		return nil, err
	}

	results, err := mapImageComponentResolverToEmbeddedImageScanComponentResolver(ctx, resolver, component, cves)
	if err != nil {
		return nil, err
	}

	return results, err
}

func (resolver *Resolver) getComponentsFromDackBox(ctx context.Context, q PaginatedQuery) ([]*EmbeddedImageScanComponentResolver, error) {
	components, err := resolver.ComponentsV2(ctx, q)
	if err != nil {
		return nil, err
	}

	ret := make([]*EmbeddedImageScanComponentResolver, 0, len(components))
	for _, component := range components {
		cves, err := resolver.getComponentCVEsFromDackBox(ctx, component.data.GetId())
		if err != nil {
			return nil, err
		}

		embedded, err := mapImageComponentResolverToEmbeddedImageScanComponentResolver(ctx, resolver, component, cves)
		if err != nil {
			return nil, err
		}
		ret = append(ret, embedded)
	}
	return ret, nil
}

func mapImageComponentResolverToEmbeddedImageScanComponentResolver(ctx context.Context, root *Resolver, component *imageComponentResolver, cves []*cVEResolver) (*EmbeddedImageScanComponentResolver, error) {
	embedded := &EmbeddedImageScanComponentResolver{
		root: root,
		data: imageComponentConverter.ProtoImageComponentToEmbeddedImageScanComponent(component.data),
	}

	embedded.updateVulns(cves)

	err := embedded.updateLastScanned(ctx, component)
	if err != nil {
		return nil, err
	}

	err = embedded.updateLayerIndex(ctx)
	if err != nil {
		return nil, err
	}
	return embedded, nil
}

func (resolver *Resolver) getComponentCVEsFromDackBox(ctx context.Context, componentID string) ([]*cVEResolver, error) {
	cveQuery := search.NewQueryBuilder().AddExactMatches(search.ComponentID, componentID).Query()
	cves, err := resolver.Cves(ctx, PaginatedQuery{Query: &cveQuery})
	if err != nil {
		return nil, err
	}
	return cves, nil
}

func (eicr *EmbeddedImageScanComponentResolver) updateLastScanned(ctx context.Context, component *imageComponentResolver) error {
	ls, err := component.getLastScannedTime(ctx)
	if err != nil {
		return err
	}

	eicr.lastScanned = ls
	return nil
}

func (eicr *EmbeddedImageScanComponentResolver) updateVulns(cves []*cVEResolver) {
	embeddedVulns := make([]*storage.EmbeddedVulnerability, 0, len(cves))
	for _, cve := range cves {
		embeddedVulns = append(embeddedVulns, converter.ProtoCVEToEmbeddedCVE(cve.data))
	}
	eicr.data.Vulns = embeddedVulns
}

func (eicr *EmbeddedImageScanComponentResolver) updateLayerIndex(ctx context.Context) error {
	cID := &componentID{
		Name:    eicr.data.GetName(),
		Version: eicr.data.GetVersion(),
	}
	// TODO: cannot figure out this part
	_, err := eicr.root.ImageComponentEdgeDataStore.SearchRawEdges(ctx, search.NewQueryBuilder().AddExactMatches(search.ComponentID, cID.toString()).ProtoQuery())
	if err != nil {
		return err
	}

	return nil
}
