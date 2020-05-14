package service

import (
	"context"

	"github.com/gogo/protobuf/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stackrox/rox/central/role/utils"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/grpc/authn"
	"github.com/stackrox/rox/pkg/grpc/authz/allow"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ClusterService is the struct that manages the cluster API
type serviceImpl struct{}

// RegisterServiceServer registers this service with the given gRPC Server.
func (s *serviceImpl) RegisterServiceServer(grpcServer *grpc.Server) {
	v1.RegisterAuthServiceServer(grpcServer, s)
}

// RegisterServiceHandler registers this service with the given gRPC Gateway endpoint.
func (s *serviceImpl) RegisterServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return v1.RegisterAuthServiceHandler(ctx, mux, conn)
}

// AuthFuncOverride specifies the auth criteria for this API.
func (s *serviceImpl) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, allow.Anonymous().Authorized(ctx, fullMethodName)
}

// GetAuthStatus retrieves the auth status based on the credentials given to the server.
func (s *serviceImpl) GetAuthStatus(ctx context.Context, request *v1.Empty) (*v1.AuthStatus, error) {
	id := authn.IdentityFromContext(ctx)
	if id == nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}

	return authStatusForID(id)
}

func convertAttributes(attrMap map[string][]string) []*v1.UserAttribute {
	if attrMap == nil {
		return nil
	}
	result := make([]*v1.UserAttribute, 0, len(attrMap))

	for k, vs := range attrMap {
		attr := &v1.UserAttribute{
			Key:    k,
			Values: vs,
		}
		result = append(result, attr)
	}
	return result
}

func fillAccessLists(userInfo *storage.UserInfo) {
	if userInfo == nil {
		return
	}
	utils.FillAccessList(userInfo.Permissions)
	for _, role := range userInfo.Roles {
		utils.FillAccessList(role)
	}
}

func authStatusForID(id authn.Identity) (*v1.AuthStatus, error) {
	exp, err := types.TimestampProto(id.Expiry())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "expiration time: %s", err)
	}

	result := &v1.AuthStatus{
		Expires:        exp,
		UserInfo:       id.User().Clone(),
		UserAttributes: convertAttributes(id.Attributes()),
	}
	fillAccessLists(result.UserInfo)

	if provider := id.ExternalAuthProvider(); provider != nil {
		// every Identity should now have an auth provider but API token Identities won't have a Backend
		if backend := provider.Backend(); backend != nil {
			result.RefreshUrl = backend.RefreshURL()
		}
		authProvider := provider.StorageView().Clone()
		if authProvider != nil {
			// config might contain semi-sensitive values, so strip it
			authProvider.Config = nil
		}
		result.AuthProvider = authProvider
	}
	if svc := id.Service(); svc != nil {
		result.Id = &v1.AuthStatus_ServiceId{ServiceId: svc}
	} else {
		result.Id = &v1.AuthStatus_UserId{UserId: id.UID()}
	}
	return result, nil
}
