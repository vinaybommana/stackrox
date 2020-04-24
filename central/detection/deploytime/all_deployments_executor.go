package deploytime

import (
	"context"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/deployment/datastore"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/detection"
	"github.com/stackrox/rox/pkg/detection/deploytime"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/utils"
)

func newAllDeploymentsExecutor(executorCtx context.Context, deployments datastore.DataStore) deploytime.AlertCollectingExecutor {
	return &allDeploymentsExecutor{
		deployments: deployments,
		executorCtx: executorCtx,
	}
}

type allDeploymentsExecutor struct {
	executorCtx context.Context
	deployments datastore.DataStore
	alerts      []*storage.Alert
}

func (d *allDeploymentsExecutor) GetAlerts() []*storage.Alert {
	if features.BooleanPolicyLogic.Enabled() {
		utils.Should(errors.New("search-based policy evaluation is deprecated"))
		return nil
	}
	return d.alerts
}

func (d *allDeploymentsExecutor) ClearAlerts() {
	if features.BooleanPolicyLogic.Enabled() {
		utils.Should(errors.New("search-based policy evaluation is deprecated"))
	}
	d.alerts = nil
}

func (d *allDeploymentsExecutor) Execute(compiled detection.CompiledPolicy) error {
	if features.BooleanPolicyLogic.Enabled() {
		return utils.Should(errors.New("search-based policy evaluation is deprecated"))
	}
	if compiled.Policy().GetDisabled() {
		return nil
	}
	violationsByDeployment, err := compiled.Matcher().Match(d.executorCtx, d.deployments)
	if err != nil {
		return err
	}
	for deploymentID, violations := range violationsByDeployment {
		dep, exists, err := d.deployments.GetDeployment(d.executorCtx, deploymentID)
		if err != nil {
			return err
		}
		if !exists {
			log.Errorf("deployment with id %q had violations, but doesn't exist", deploymentID)
			continue
		}
		if !compiled.AppliesTo(dep) {
			continue
		}
		d.alerts = append(d.alerts, deploytime.PolicyDeploymentAndViolationsToAlert(compiled.Policy(), dep, violations.AlertViolations))
	}
	return nil
}
