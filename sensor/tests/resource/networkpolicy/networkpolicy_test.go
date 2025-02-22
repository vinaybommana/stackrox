package networkpolicy

import (
	"log"
	"testing"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/sensor/tests/helper"
	"github.com/stackrox/rox/sensor/testutils"
	"github.com/stretchr/testify/suite"
	"sigs.k8s.io/e2e-framework/klient/k8s"
)

var (
	NginxDeployment            = helper.K8sResourceInfo{Kind: "Deployment", YamlFile: "nginx.yaml"}
	IngressPolicyAllow443      = helper.K8sResourceInfo{Kind: "NetworkPolicy", YamlFile: "netpol-allow-443.yaml"}
	EgressPolicyBlockAllEgress = helper.K8sResourceInfo{Kind: "NetworkPolicy", YamlFile: "netpol-block-egress.yaml"}
)

type NetworkPolicySuite struct {
	testContext *helper.TestContext
	suite.Suite
}

func Test_NetworkPolicy(t *testing.T) {
	suite.Run(t, new(NetworkPolicySuite))
}

var _ suite.SetupAllSuite = &NetworkPolicySuite{}
var _ suite.TearDownTestSuite = &NetworkPolicySuite{}

func (s *NetworkPolicySuite) SetupSuite() {
	s.T().Setenv("ROX_RESYNC_DISABLED", "true")
	policies, err := testutils.GetPoliciesFromFile("data/policies.json")
	if err != nil {
		log.Fatalln(err)
	}
	cfg := helper.DefaultCentralConfig()
	cfg.InitialSystemPolicies = policies

	if testContext, err := helper.NewContextWithConfig(s.T(), cfg); err != nil {
		s.Fail("failed to setup test context: %s", err)
	} else {
		s.testContext = testContext
	}
}

func (s *NetworkPolicySuite) TearDownTest() {
	s.testContext.GetFakeCentral().ClearReceivedBuffer()
}

var (
	ingressNetpolViolationName = "Deployments should have at least one ingress Network Policy"
	egressNetpolViolationName  = "Deployments should have at least one egress Network Policy"
)

func checkViolations(violations []string) func(result *central.AlertResults) error {
	return func(result *central.AlertResults) error {
		missing := set.NewStringSet(violations...)
		for _, alertMessage := range result.GetAlerts() {
			missing.Remove(alertMessage.GetPolicy().GetName())
		}

		if len(missing) != 0 {
			return errors.Errorf("expected violations not found: %v", missing.AsSlice())
		}
		return nil
	}
}

func (s *NetworkPolicySuite) Test_Deployment_NetpolViolations() {
	testCases := map[string]struct {
		netpolsApplied     []helper.K8sResourceInfo
		violationsExpected []string
	}{
		"No policies applied: should have two violations": {
			netpolsApplied:     []helper.K8sResourceInfo{},
			violationsExpected: []string{ingressNetpolViolationName, egressNetpolViolationName},
		},
		"Both policies applied: should have no violations": {
			netpolsApplied:     []helper.K8sResourceInfo{IngressPolicyAllow443, EgressPolicyBlockAllEgress},
			violationsExpected: []string{},
		},
		"Ingress applied: egress violation": {
			netpolsApplied:     []helper.K8sResourceInfo{IngressPolicyAllow443},
			violationsExpected: []string{egressNetpolViolationName},
		},
		"Egress applied: ingress violation": {
			netpolsApplied:     []helper.K8sResourceInfo{EgressPolicyBlockAllEgress},
			violationsExpected: []string{ingressNetpolViolationName},
		},
	}

	for name, testCase := range testCases {
		s.T().Run(name, func(t *testing.T) {
			resourcesToApply := append(testCase.netpolsApplied, NginxDeployment)
			s.testContext.RunTest(
				helper.WithResources(resourcesToApply),
				helper.WithTestCase(func(t *testing.T, tc *helper.TestContext, objects map[string]k8s.Object) {
					tc.LastViolationState("nginx-deployment", checkViolations(testCase.violationsExpected), name)
				}))
		})

	}
}
