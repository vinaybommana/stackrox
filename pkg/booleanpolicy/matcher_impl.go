package booleanpolicy

import (
	"context"
	"fmt"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/predicate"
	"github.com/stackrox/rox/pkg/searchbasedpolicies"
	"github.com/stackrox/rox/pkg/searchbasedpolicies/builders"
)

type predicateTuple struct {
	processPredicate    predicate.Predicate
	deploymentPredicate predicate.Predicate
	imagePredicate      predicate.Predicate
}

type matcherImpl struct {
	predicates []predicateTuple
}

func matchWithPredicateTuple(predicates predicateTuple, deployment *storage.Deployment, images []*storage.Image, indicator *storage.ProcessIndicator) ([]*storage.Alert_Violation, error) {
	var results []*search.Result
	if indicator != nil {
		result, matches := predicates.processPredicate.Evaluate(indicator)
		if !matches {
			return nil, nil
		}
		results = append(results, result)
	}

	if deployment != nil {
		result, matches := predicates.deploymentPredicate.Evaluate(deployment)
		if !matches {
			return nil, nil
		}
		results = append(results, result)
	}

	if len(images) > 0 {
		var foundMatch bool
		for _, img := range images {
			result, matches := predicates.imagePredicate.Evaluate(img)
			if matches {
				foundMatch = true
				results = append(results, result)
			}
		}
		if !foundMatch {
			return nil, nil
		}
	}

	finalResult := predicate.MergeResults(results...)

	// TODO(viswa): Figure out how to populate these for Boolean policies.
	violations := []*storage.Alert_Violation{{Message: fmt.Sprintf("TODO (%+v)", finalResult)}}
	return violations, nil

}

// MatchOne returns detection against the deployment and images using predicate matching
// The deployment parameter can be nil in the case of image detection
func (m *matcherImpl) MatchOne(ctx context.Context, deployment *storage.Deployment, images []*storage.Image, indicator *storage.ProcessIndicator) (searchbasedpolicies.Violations, error) {
	var allViolations []*storage.Alert_Violation
	for _, tuple := range m.predicates {
		violations, err := matchWithPredicateTuple(tuple, deployment, images, indicator)
		if err != nil {
			return searchbasedpolicies.Violations{}, err
		}
		allViolations = append(allViolations, violations...)
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
