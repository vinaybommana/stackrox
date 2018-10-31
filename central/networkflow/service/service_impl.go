package service

import (
	"context"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	deploymentDataStore "github.com/stackrox/rox/central/deployment/datastore"
	networkFlowStore "github.com/stackrox/rox/central/networkflow/store"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/grpc/authz"
	"github.com/stackrox/rox/pkg/grpc/authz/perrpc"
	"github.com/stackrox/rox/pkg/grpc/authz/user"
	"github.com/stackrox/rox/pkg/protoconv"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	authorizer = perrpc.FromMap(map[authz.Authorizer][]string{
		user.With(permissions.View(resources.NetworkGraph)): {
			"/v1.NetworkGraphService/GetNetworkGraph",
		},
	})
	defaultSince = -5 * time.Minute
)

// serviceImpl provides APIs for alerts.
type serviceImpl struct {
	clusterStore networkFlowStore.ClusterStore
	deployments  deploymentDataStore.DataStore
}

// RegisterServiceServer registers this service with the given gRPC Server.
func (s *serviceImpl) RegisterServiceServer(grpcServer *grpc.Server) {
	v1.RegisterNetworkGraphServiceServer(grpcServer, s)
}

// RegisterServiceHandler registers this service with the given gRPC Gateway endpoint.
func (s *serviceImpl) RegisterServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return v1.RegisterNetworkGraphServiceHandler(ctx, mux, conn)
}

// AuthFuncOverride specifies the auth criteria for this API.
func (s *serviceImpl) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, authorizer.Authorized(ctx, fullMethodName)
}

func (s *serviceImpl) GetNetworkGraph(context context.Context, request *v1.NetworkGraphRequest) (*v1.NetworkGraph, error) {
	if request.GetClusterId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "cluster ID must be specified")
	}

	since := timestamp.FromProtobuf(request.GetSince())
	if since == 0 {
		since = timestamp.FromGoTime(time.Now().Add(defaultSince))
	}
	// Get the deployments we want to check connectivity between.
	deployments, err := s.getDeployments(request.GetClusterId(), request.GetQuery())

	if err != nil {
		return nil, err
	}

	// compute nodes
	var nodes []*v1.NetworkNode

	nodeIndices := make(map[string]int)
	for i, d := range deployments {
		nodes = append(nodes, &v1.NetworkNode{
			DeploymentId:   d.GetId(),
			DeploymentName: d.GetName(),
			Cluster:        d.GetClusterName(),
			Namespace:      d.GetNamespace(),
			OutEdges:       make(map[int32]*v1.NetworkEdgePropertiesBundle),
		})
		nodeIndices[d.GetId()] = i
	}

	flowStore := s.clusterStore.GetFlowStore(request.GetClusterId())

	if flowStore == nil {
		return nil, status.Errorf(codes.NotFound, "no flows found for cluster %s", request.GetClusterId())
	}

	flows, _, err := flowStore.GetAllFlows()
	if err != nil {
		return nil, err
	}

	// compute edges

	// Filter by deployments, and then by time.
	filteredFlows := filterNetworkFlowsByDeployments(flows, deployments)
	filteredFlows = filterNetworkFlowsByTime(filteredFlows, since)

	for _, flow := range filteredFlows {
		props := flow.GetProps()
		srcIdx := nodeIndices[props.GetSrcDeploymentId()]
		srcNode := nodes[srcIdx]
		tgtIdx := int32(nodeIndices[props.GetDstDeploymentId()])

		tgtEdgeBundle := srcNode.OutEdges[tgtIdx]
		if tgtEdgeBundle == nil {
			tgtEdgeBundle = &v1.NetworkEdgePropertiesBundle{}
			srcNode.OutEdges[tgtIdx] = tgtEdgeBundle
		}

		edgeProps := &v1.NetworkEdgeProperties{
			Port:     props.GetDstPort(),
			Protocol: props.L4Protocol,
		}

		edgeProps.LastActiveTimestamp = flow.GetLastSeenTimestamp()
		if edgeProps.LastActiveTimestamp == nil {
			edgeProps.LastActiveTimestamp = protoconv.ConvertTimeToTimestamp(time.Now())
		}

		tgtEdgeBundle.Properties = append(tgtEdgeBundle.Properties, edgeProps)
	}

	return &v1.NetworkGraph{
		Nodes: nodes,
	}, nil
}

func filterNetworkFlowsByDeployments(flows []*v1.NetworkFlow, deployments []*v1.Deployment) (filtered []*v1.NetworkFlow) {

	filtered = flows[:0]
	deploymentIDMap := make(map[string]bool)
	for _, d := range deployments {
		deploymentIDMap[d.Id] = true
	}

	for _, flow := range flows {
		srcID := flow.GetProps().SrcDeploymentId
		dstID := flow.GetProps().DstDeploymentId

		if deploymentIDMap[srcID] && deploymentIDMap[dstID] {
			filtered = append(filtered, flow)
		}

	}

	return
}

func filterNetworkFlowsByTime(flows []*v1.NetworkFlow, since timestamp.MicroTS) (filtered []*v1.NetworkFlow) {
	filtered = flows[:0]

	for _, flow := range flows {
		flowTS := timestamp.FromProtobuf(flow.LastSeenTimestamp)
		if flowTS != 0 && flowTS < since {
			continue
		}
		filtered = append(filtered, flow)
	}

	return
}

func (s *serviceImpl) getDeployments(clusterID string, query string) (deployments []*v1.Deployment, err error) {
	clusterQuery := search.NewQueryBuilder().AddStrings(search.ClusterID, clusterID).ProtoQuery()

	q := clusterQuery
	if query != "" {
		q, err = search.ParseRawQuery(query)
		if err != nil {
			return
		}
		q = search.ConjunctionQuery(q, clusterQuery)
	}

	deployments, err = s.deployments.SearchRawDeployments(q)
	return
}
