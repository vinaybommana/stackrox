package service

import (
	"context"
	"sort"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stackrox/rox/central/compliance/manager"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/grpc/authz"
	"github.com/stackrox/rox/pkg/grpc/authz/perrpc"
	"github.com/stackrox/rox/pkg/grpc/authz/user"
	"github.com/stackrox/rox/pkg/protoconv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	authorizer = perrpc.FromMap(map[authz.Authorizer][]string{
		user.With(permissions.View(resources.ComplianceRuns)): {
			"/v1.ComplianceManagementService/GetRecentRuns",
			"/v1.ComplianceManagementService/GetRunStatuses",
		},
		user.With(permissions.Modify(resources.ComplianceRuns)): {
			"/v1.ComplianceManagementService/TriggerRun",
			"/v1.ComplianceManagementService/TriggerRuns",
		},
		user.With(permissions.View(resources.ComplianceRunSchedule)): {
			"/v1.ComplianceManagementService/GetRunSchedules",
		},
		user.With(permissions.Modify(resources.ComplianceRunSchedule)): {
			"/v1.ComplianceManagementService/AddRunSchedule",
			"/v1.ComplianceManagementService/UpdateRunSchedule",
			"/v1.ComplianceManagementService/DeleteRunSchedule",
		},
	})
)

type service struct {
	manager manager.ComplianceManager
}

func newService(manager manager.ComplianceManager) *service {
	return &service{
		manager: manager,
	}
}

func (s *service) RegisterServiceServer(server *grpc.Server) {
	v1.RegisterComplianceManagementServiceServer(server, s)
}

func (s *service) RegisterServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return v1.RegisterComplianceManagementServiceHandler(ctx, mux, conn)
}

func (s *service) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, authorizer.Authorized(ctx, fullMethodName)
}

func (s *service) AddRunSchedule(ctx context.Context, req *v1.AddComplianceRunScheduleRequest) (*v1.AddComplianceRunScheduleResponse, error) {
	schedule, err := s.manager.AddSchedule(req.GetScheduleSpec())
	if err != nil {
		return nil, err
	}
	return &v1.AddComplianceRunScheduleResponse{
		AddedSchedule: schedule,
	}, nil
}

func (s *service) UpdateRunSchedule(ctx context.Context, req *v1.UpdateComplianceRunScheduleRequest) (*v1.UpdateComplianceRunScheduleResponse, error) {
	if req.GetUpdatedSpec().GetId() == "" {
		req.UpdatedSpec.Id = req.GetScheduleId()
	} else if req.GetUpdatedSpec().GetId() != req.GetScheduleId() {
		return nil, status.Errorf(codes.InvalidArgument, "id in updated spec body must be empty or match schedule id %q, is: %q", req.GetScheduleId(), req.GetUpdatedSpec().GetId())
	}

	schedule, err := s.manager.UpdateSchedule(req.GetUpdatedSpec())
	if err != nil {
		return nil, err
	}
	return &v1.UpdateComplianceRunScheduleResponse{
		UpdatedSchedule: schedule,
	}, nil
}

func (s *service) DeleteRunSchedule(ctx context.Context, req *v1.DeleteComplianceRunScheduleRequest) (*v1.Empty, error) {
	err := s.manager.DeleteSchedule(req.GetScheduleId())
	if err != nil {
		return nil, err
	}
	return &v1.Empty{}, nil
}

func (s *service) GetRecentRuns(ctx context.Context, req *v1.GetRecentComplianceRunsRequest) (*v1.GetRecentComplianceRunsResponse, error) {
	runs := s.manager.GetRecentRuns(req)
	sort.Slice(runs, func(i, j int) bool {
		return protoconv.CompareProtoTimestamps(runs[i].StartTime, runs[j].StartTime) < 0
	})

	return &v1.GetRecentComplianceRunsResponse{
		ComplianceRuns: runs,
	}, nil
}

func (s *service) GetRunSchedules(ctx context.Context, req *v1.GetComplianceRunSchedulesRequest) (*v1.GetComplianceRunSchedulesResponse, error) {
	schedules := s.manager.GetSchedules(req)
	return &v1.GetComplianceRunSchedulesResponse{
		Schedules: schedules,
	}, nil
}

func (s *service) TriggerRun(ctx context.Context, req *v1.TriggerComplianceRunRequest) (*v1.TriggerComplianceRunResponse, error) {
	run, err := s.manager.TriggerRun(req.GetClusterId(), req.GetStandardId())
	if err != nil {
		return nil, err
	}
	return &v1.TriggerComplianceRunResponse{
		StartedRun: run,
	}, nil
}

func (s *service) TriggerRuns(ctx context.Context, req *v1.TriggerComplianceRunsRequest) (*v1.TriggerComplianceRunsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *service) GetRunStatuses(ctx context.Context, req *v1.GetComplianceRunStatusesRequest) (*v1.GetComplianceRunStatusesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
