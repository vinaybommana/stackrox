package booleanpolicy

import (
	"context"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/booleanpolicy/augmentedobjs"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/pathutil"
	"github.com/stackrox/rox/pkg/booleanpolicy/violations"
	"github.com/stackrox/rox/pkg/searchbasedpolicies"
	"github.com/stackrox/rox/pkg/searchbasedpolicies/builders"
)

type matcherImpl struct {
	evaluators []sectionAndEvaluator
	stage      storage.LifecycleStage
}

func matchWithEvaluator(sectionAndEval sectionAndEvaluator, obj *pathutil.AugmentedObj, stage storage.LifecycleStage) ([]*storage.Alert_Violation, error) {
	finalResult, matched := sectionAndEval.evaluator.Evaluate(obj.Value())
	if !matched {
		return nil, nil
	}
	return violations.ViolationPrinter(stage, sectionAndEval.sectionName, finalResult)
}

func (m *matcherImpl) MatchImage(_ context.Context, image *storage.Image) (searchbasedpolicies.Violations, error) {
	obj, err := augmentedobjs.ConstructImage(image)
	if err != nil {
		return searchbasedpolicies.Violations{}, err
	}
	violations, err := m.getViolations(obj)
	if err != nil || violations == nil {
		return searchbasedpolicies.Violations{}, err
	}
	return *violations, nil
}

func (m *matcherImpl) getViolations(obj *pathutil.AugmentedObj) (*searchbasedpolicies.Violations, error) {
	var allViolations []*storage.Alert_Violation
	var atLeastOneMatched bool
	for _, eval := range m.evaluators {
		violations, err := matchWithEvaluator(eval, obj, m.stage)
		if err != nil {
			return nil, err
		}
		atLeastOneMatched = atLeastOneMatched || len(violations) > 0
		allViolations = append(allViolations, violations...)
	}
	if !atLeastOneMatched {
		return nil, nil
	}
	return &searchbasedpolicies.Violations{
		AlertViolations: allViolations,
	}, nil
}

func (m *matcherImpl) MatchDeploymentWithProcess(_ context.Context, deployment *storage.Deployment, images []*storage.Image, indicator *storage.ProcessIndicator, processOutsideWhitelist bool) (searchbasedpolicies.Violations, error) {
	obj, err := augmentedobjs.ConstructDeploymentWithProcess(deployment, images, indicator, processOutsideWhitelist)
	if err != nil {
		return searchbasedpolicies.Violations{}, err
	}

	violations, err := m.getViolations(obj)
	if err != nil || violations == nil {
		return searchbasedpolicies.Violations{}, err
	}
	v := &storage.Alert_ProcessViolation{Processes: []*storage.ProcessIndicator{indicator}}
	builders.UpdateRuntimeAlertViolationMessage(v)
	violations.ProcessViolation = v
	return *violations, nil
}

// MatchDeployment runs detection against the deployment and images.
func (m *matcherImpl) MatchDeployment(_ context.Context, deployment *storage.Deployment, images []*storage.Image) (searchbasedpolicies.Violations, error) {
	obj, err := augmentedobjs.ConstructDeployment(deployment, images)
	if err != nil {
		return searchbasedpolicies.Violations{}, err
	}
	violations, err := m.getViolations(obj)
	if err != nil || violations == nil {
		return searchbasedpolicies.Violations{}, err
	}
	return *violations, nil
}
