package booleanpolicy

import (
	"context"
	"fmt"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/booleanpolicy/augmentedobjs"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/pathutil"
	"github.com/stackrox/rox/pkg/searchbasedpolicies"
	"github.com/stackrox/rox/pkg/searchbasedpolicies/builders"
)

type matcherImpl struct {
	evaluators []sectionAndEvaluator
}

func matchWithEvaluator(sectionAndEval sectionAndEvaluator, obj *pathutil.AugmentedObj) ([]*storage.Alert_Violation, error) {
	finalResult, matched := sectionAndEval.evaluator.Evaluate(obj.Value())
	if !matched {
		return nil, nil
	}

	// TODO(viswa): Figure out how to populate these for Boolean policies.
	violations := []*storage.Alert_Violation{{Message: fmt.Sprintf("TODO (%+v)", finalResult)}}
	return violations, nil

}

func (m *matcherImpl) MatchImage(ctx context.Context, image *storage.Image) (searchbasedpolicies.Violations, error) {
	var allViolations []*storage.Alert_Violation
	obj, err := augmentedobjs.ConstructImage(image)
	if err != nil {
		return searchbasedpolicies.Violations{}, err
	}
	for _, eval := range m.evaluators {
		result, matched := eval.evaluator.Evaluate(obj.Value())
		if matched {
			violations := []*storage.Alert_Violation{{Message: fmt.Sprintf("TODO (%+v)", result)}}
			allViolations = append(allViolations, violations...)
		}
	}
	// The following line automatically handles the case where there is no match, since allViolations will be nil.
	return searchbasedpolicies.Violations{AlertViolations: allViolations}, nil
}

// MatchOne returns detection against the deployment and images using predicate matching
// The deployment parameter can be nil in the case of image detection
func (m *matcherImpl) MatchDeployment(ctx context.Context, deployment *storage.Deployment, images []*storage.Image, indicator *storage.ProcessIndicator) (searchbasedpolicies.Violations, error) {
	var allViolations []*storage.Alert_Violation
	var atLeastOneMatched bool
	obj, err := augmentedobjs.ConstructDeployment(deployment, images, indicator)
	if err != nil {
		return searchbasedpolicies.Violations{}, err
	}

	for _, eval := range m.evaluators {
		violations, err := matchWithEvaluator(eval, obj)
		if err != nil {
			return searchbasedpolicies.Violations{}, err
		}
		atLeastOneMatched = atLeastOneMatched || len(violations) > 0
		allViolations = append(allViolations, violations...)
	}
	if !atLeastOneMatched {
		return searchbasedpolicies.Violations{}, nil
	}
	violations := searchbasedpolicies.Violations{
		AlertViolations: allViolations,
	}
	if indicator != nil {
		v := &storage.Alert_ProcessViolation{Processes: []*storage.ProcessIndicator{indicator}}
		builders.UpdateRuntimeAlertViolationMessage(v)
		violations.ProcessViolation = v
	}
	return violations, nil
}
