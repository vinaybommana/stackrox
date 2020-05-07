package booleanpolicy

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/pathutil"
	"github.com/stackrox/rox/pkg/searchbasedpolicies"
	"github.com/stackrox/rox/pkg/searchbasedpolicies/builders"
)

type matcherImpl struct {
	evaluators []sectionAndEvaluator
}

func findMatchingContainerIdxForIndicator(deployment *storage.Deployment, indicator *storage.ProcessIndicator) (int, error) {
	for i, container := range deployment.GetContainers() {
		if container.GetName() == indicator.GetContainerName() {
			return i, nil
		}
	}
	return 0, errors.Errorf("indicator %s could not be matched (container name %s not found in deployment %s/%s/%s",
		indicator.GetSignal().GetExecFilePath(), indicator.GetContainerName(), deployment.GetClusterId(), deployment.GetNamespace(), deployment.GetName())

}

func matchWithEvaluator(sectionAndEval sectionAndEvaluator, deployment *storage.Deployment, images []*storage.Image, indicator *storage.ProcessIndicator) ([]*storage.Alert_Violation, error) {
	obj := pathutil.NewAugmentedObj(deployment)
	if len(images) != len(deployment.GetContainers()) {
		return nil, errors.Errorf("deployment %s/%s had %d containers, but got %d images",
			deployment.GetNamespace(), deployment.GetName(), len(deployment.GetContainers()), len(images))
	}
	for i, image := range images {
		err := obj.AddPlainObjAt(
			(&pathutil.Path{}).TraverseField("Containers").IndexSlice(i).TraverseField(imageAugmentKey),
			image)
		if err != nil {
			return nil, err
		}
	}

	if indicator != nil {
		matchingContainerIdx, err := findMatchingContainerIdxForIndicator(deployment, indicator)
		if err != nil {
			return nil, err
		}
		err = obj.AddPlainObjAt(
			(&pathutil.Path{}).TraverseField("Containers").IndexSlice(matchingContainerIdx).TraverseField(processAugmentKey),
			indicator,
		)
		if err != nil {
			return nil, err
		}
	}

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
	for _, eval := range m.evaluators {
		result, matched := eval.evaluator.Evaluate(pathutil.NewAugmentedObj(image).Value())
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
	for _, eval := range m.evaluators {
		violations, err := matchWithEvaluator(eval, deployment, images, indicator)
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
