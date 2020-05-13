package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	lifecycleMocks "github.com/stackrox/rox/central/detection/lifecycle/mocks"
	"github.com/stackrox/rox/central/policy/datastore/mocks"
	connectionMocks "github.com/stackrox/rox/central/sensor/service/connection/mocks"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/booleanpolicy"
	detectionMocks "github.com/stackrox/rox/pkg/detection/mocks"
	"github.com/stackrox/rox/pkg/features"
	matcherMocks "github.com/stackrox/rox/pkg/searchbasedpolicies/matcher/mocks"
	sbpMocks "github.com/stackrox/rox/pkg/searchbasedpolicies/mocks"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/status"
)

var (
	mockRequestOneID = &v1.ExportPoliciesRequest{
		PolicyIds: []string{"Joseph Rules"},
	}
	mockRequestTwoIDs = &v1.ExportPoliciesRequest{
		PolicyIds: []string{"Joseph Rules", "abcd"},
	}
)

func TestPolicyService(t *testing.T) {
	suite.Run(t, new(PolicyServiceTestSuite))
}

type testDeploymentMatcher struct {
	*detectionMocks.MockPolicySet
}

func (t *testDeploymentMatcher) RemoveNotifier(_ string) error {
	return nil
}

type PolicyServiceTestSuite struct {
	suite.Suite
	policies                     *mocks.MockDataStore
	mockDeploymentMatcherBuilder *matcherMocks.MockBuilder
	mockBuildTimePolicies        *detectionMocks.MockPolicySet
	mockLifecycleManager         *lifecycleMocks.MockManager
	mockConnectionManager        *connectionMocks.MockManager
	tested                       Service

	envIsolator *testutils.EnvIsolator

	mockCtrl *gomock.Controller
}

func (s *PolicyServiceTestSuite) SetupTest() {
	s.envIsolator = testutils.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(features.PolicyImportExport.EnvVar(), "true")

	s.mockCtrl = gomock.NewController(s.T())

	s.policies = mocks.NewMockDataStore(s.mockCtrl)

	s.mockDeploymentMatcherBuilder = matcherMocks.NewMockBuilder(s.mockCtrl)
	s.mockBuildTimePolicies = detectionMocks.NewMockPolicySet(s.mockCtrl)
	s.mockLifecycleManager = lifecycleMocks.NewMockManager(s.mockCtrl)
	s.mockConnectionManager = connectionMocks.NewMockManager(s.mockCtrl)

	s.tested = New(
		s.policies,
		nil,
		nil,
		nil,
		nil,
		&testDeploymentMatcher{s.mockBuildTimePolicies},
		s.mockDeploymentMatcherBuilder,
		nil,
		s.mockLifecycleManager,
		nil,
		nil,
		nil,
		s.mockConnectionManager,
	)
}

func (s *PolicyServiceTestSuite) TearDownTest() {
	defer s.envIsolator.RestoreAll()
	s.mockCtrl.Finish()
}

func (s *PolicyServiceTestSuite) compareErrorsToExpected(expectedErrors []*v1.ExportPolicyError, apiError error) {
	apiStatus, ok := status.FromError(apiError)
	s.Require().True(ok)
	details := apiStatus.Details()
	s.Len(details, 1)
	exportErrors, ok := details[0].(*v1.ExportPoliciesErrorList)
	s.Require().True(ok)
	// actual errors == expected errors ignoring order
	s.Len(exportErrors.GetErrors(), len(expectedErrors))
	for _, expected := range expectedErrors {
		s.Contains(exportErrors.GetErrors(), expected)
	}
}

func makeError(errorID, errorString string) *v1.ExportPolicyError {
	return &v1.ExportPolicyError{
		PolicyId: errorID,
		Error: &v1.PolicyError{
			Error: errorString,
		},
	}
}

func (s *PolicyServiceTestSuite) TestExportInvalidIDFails() {
	ctx := context.Background()
	mockErrors := []*v1.ExportPolicyError{
		makeError(mockRequestOneID.PolicyIds[0], "not found"),
	}
	s.policies.EXPECT().GetPolicies(ctx, mockRequestOneID.PolicyIds).Return(make([]*storage.Policy, 0), []int{0}, []error{errors.New("not found")}, nil)
	resp, err := s.tested.ExportPolicies(ctx, mockRequestOneID)
	s.Nil(resp)
	s.Error(err)
	s.compareErrorsToExpected(mockErrors, err)
}

func (s *PolicyServiceTestSuite) TestExportValidIDSucceeds() {
	ctx := context.Background()
	mockPolicy := &storage.Policy{
		Id: mockRequestOneID.PolicyIds[0],
	}
	s.policies.EXPECT().GetPolicies(ctx, mockRequestOneID.PolicyIds).Return([]*storage.Policy{mockPolicy}, nil, nil, nil)
	resp, err := s.tested.ExportPolicies(ctx, mockRequestOneID)
	s.NoError(err)
	s.NotNil(resp)
	s.Len(resp.GetPolicies(), 1)
	s.Equal(mockPolicy, resp.Policies[0])
}

func (s *PolicyServiceTestSuite) TestExportMixedSuccessAndMissing() {
	ctx := context.Background()
	mockPolicy := &storage.Policy{
		Id: mockRequestTwoIDs.PolicyIds[0],
	}
	mockErrors := []*v1.ExportPolicyError{
		makeError(mockRequestTwoIDs.PolicyIds[1], "not found"),
	}
	s.policies.EXPECT().GetPolicies(ctx, mockRequestTwoIDs.PolicyIds).Return([]*storage.Policy{mockPolicy}, []int{1}, []error{errors.New("not found")}, nil)
	resp, err := s.tested.ExportPolicies(ctx, mockRequestTwoIDs)
	s.Nil(resp)
	s.Error(err)
	s.compareErrorsToExpected(mockErrors, err)
}

func (s *PolicyServiceTestSuite) TestExportMultipleFailures() {
	ctx := context.Background()
	errString := "test"
	storeErrors := []error{errors.New(errString), errors.New("not found")}
	mockErrors := []*v1.ExportPolicyError{
		makeError(mockRequestTwoIDs.PolicyIds[0], errString),
		makeError(mockRequestTwoIDs.PolicyIds[1], "not found"),
	}
	s.policies.EXPECT().GetPolicies(ctx, mockRequestTwoIDs.PolicyIds).Return(make([]*storage.Policy, 0), []int{0, 1}, storeErrors, nil)
	resp, err := s.tested.ExportPolicies(ctx, mockRequestTwoIDs)
	s.Nil(resp)
	s.Error(err)
	s.compareErrorsToExpected(mockErrors, err)
}

func (s *PolicyServiceTestSuite) TestDryRunRuntime() {
	ctx := context.Background()
	runtimePolicy := &storage.Policy{
		Id:              "1",
		Name:            "RuntimePolicy",
		Severity:        storage.Severity_LOW_SEVERITY,
		LifecycleStages: []storage.LifecycleStage{storage.LifecycleStage_RUNTIME},
		Categories:      []string{"test"},
		Fields: &storage.PolicyFields{
			ProcessPolicy: &storage.ProcessPolicy{
				Name: "apt-get",
			},
			SetPrivileged: &storage.PolicyFields_Privileged{
				Privileged: true,
			},
		},
	}
	// Runtime policy dry run exits early and returns empty results
	s.mockDeploymentMatcherBuilder.EXPECT().ForPolicy(runtimePolicy).Return(sbpMocks.NewMockMatcher(s.mockCtrl), nil).Times(1)
	resp, err := s.tested.DryRunPolicy(ctx, runtimePolicy)
	s.Nil(err)
	s.Nil(resp.GetAlerts())
}

func (s *PolicyServiceTestSuite) TestImportPolicy() {
	envIsolator := testutils.NewEnvIsolator(s.T())
	envIsolator.Setenv(features.BooleanPolicyLogic.EnvVar(), "false")
	defer envIsolator.RestoreAll()

	mockID := "1"
	mockName := "legacy policy"
	mockSeverity := storage.Severity_LOW_SEVERITY
	mockLCStages := []storage.LifecycleStage{storage.LifecycleStage_RUNTIME}
	mockCategories := []string{"test"}
	importedPolicy := &storage.Policy{
		Id:              mockID,
		Name:            mockName,
		Severity:        mockSeverity,
		LifecycleStages: mockLCStages,
		Categories:      mockCategories,
		Fields: &storage.PolicyFields{
			ProcessPolicy: &storage.ProcessPolicy{
				Name: "apt-get",
			},
			SetPrivileged: &storage.PolicyFields_Privileged{
				Privileged: true,
			},
		},
	}

	ctx := context.Background()
	mockImportResp := []*v1.ImportPolicyResponse{
		{
			Succeeded: true,
			Policy:    importedPolicy,
			Errors:    nil,
		},
	}

	s.mockDeploymentMatcherBuilder.EXPECT().ForPolicy(importedPolicy).Return(sbpMocks.NewMockMatcher(s.mockCtrl), nil).Times(1)
	s.policies.EXPECT().ImportPolicies(ctx, []*storage.Policy{importedPolicy}, false).Return(mockImportResp, true, nil)
	s.mockBuildTimePolicies.EXPECT().RemovePolicy(importedPolicy.GetId()).Return(nil)
	s.mockLifecycleManager.EXPECT().UpsertPolicy(importedPolicy).Return(nil)
	s.policies.EXPECT().GetAllPolicies(gomock.Any()).Return(nil, nil)
	s.mockConnectionManager.EXPECT().BroadcastMessage(gomock.Any())
	resp, err := s.tested.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{importedPolicy},
	})
	s.NoError(err)
	s.True(resp.AllSucceeded)
	s.Require().Len(resp.GetResponses(), 1)
	policyResp := resp.GetResponses()[0]
	resultPolicy := policyResp.GetPolicy()
	s.Equal(importedPolicy.GetFields(), resultPolicy.GetFields())
	s.Equal(importedPolicy.GetPolicySections(), resultPolicy.GetPolicySections())
}

func (s *PolicyServiceTestSuite) TestImportAndUpgradePolicy() {
	envIsolator := testutils.NewEnvIsolator(s.T())
	envIsolator.Setenv(features.BooleanPolicyLogic.EnvVar(), "true")
	defer envIsolator.RestoreAll()

	mockID := "1"
	mockName := "legacy policy"
	mockSeverity := storage.Severity_LOW_SEVERITY
	mockLCStages := []storage.LifecycleStage{storage.LifecycleStage_RUNTIME}
	mockCategories := []string{"test"}

	importedPolicy := &storage.Policy{
		Id:              mockID,
		Name:            mockName,
		Severity:        mockSeverity,
		LifecycleStages: mockLCStages,
		Categories:      mockCategories,
		Fields: &storage.PolicyFields{
			ProcessPolicy: &storage.ProcessPolicy{
				Name: "apt-get",
			},
			SetPrivileged: &storage.PolicyFields_Privileged{
				Privileged: true,
			},
		},
	}
	importRespPolicy := &storage.Policy{
		Id:              mockID,
		Name:            mockName,
		Severity:        mockSeverity,
		LifecycleStages: mockLCStages,
		Categories:      mockCategories,
		PolicyVersion:   booleanpolicy.Version,
		PolicySections: []*storage.PolicySection{
			{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Privileged",
						Values: []*storage.PolicyValue{
							{
								Value: "true",
							},
						},
					},
					{
						FieldName: "Process Name",
						Values: []*storage.PolicyValue{
							{
								Value: "apt-get",
							},
						},
					},
				},
			},
		},
	}
	ctx := context.Background()
	mockImportResp := []*v1.ImportPolicyResponse{
		{
			Succeeded: true,
			Policy:    importRespPolicy,
			Errors:    nil,
		},
	}

	s.policies.EXPECT().ImportPolicies(ctx, []*storage.Policy{importRespPolicy}, false).Return(mockImportResp, true, nil)
	s.mockBuildTimePolicies.EXPECT().RemovePolicy(importRespPolicy.GetId()).Return(nil)
	s.mockLifecycleManager.EXPECT().UpsertPolicy(importRespPolicy).Return(nil)
	s.policies.EXPECT().GetAllPolicies(gomock.Any()).Return(nil, nil)
	s.mockConnectionManager.EXPECT().BroadcastMessage(gomock.Any())
	resp, err := s.tested.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{importedPolicy},
	})
	s.NoError(err)
	s.True(resp.AllSucceeded)
	s.Require().Len(resp.GetResponses(), 1)
	policyResp := resp.GetResponses()[0]
	resultPolicy := policyResp.GetPolicy()
	s.Equal(importRespPolicy.GetFields(), resultPolicy.GetFields())
	s.Equal(importRespPolicy.GetPolicySections(), resultPolicy.GetPolicySections())
}
