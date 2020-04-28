package booleanpolicy

import (
	"context"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search/predicate"
	"github.com/stackrox/rox/pkg/searchbasedpolicies"
)

var (
	imageFactory      = predicate.NewFactory("image", (*storage.Image)(nil))
	deploymentFactory = predicate.NewFactory("deployment", (*storage.Deployment)(nil))
	processFactory    = predicate.NewFactory("process_indicator", (*storage.ProcessIndicator)(nil))
)

// Matcher matches objects against a policy.
//go:generate mockgen-wrapper
type Matcher interface {
	// MatchOne matches the policy against the passed deployment and images
	MatchOne(ctx context.Context, deployment *storage.Deployment, images []*storage.Image, pi *storage.ProcessIndicator) (searchbasedpolicies.Violations, error)
}

// BuildMatcher builds a matcher for the given policy.
// It returns an error if the policy is ill-formed.
func BuildMatcher(p *storage.Policy) (Matcher, error) {
	if err := Validate(p); err != nil {
		return nil, err
	}

	queries, err := policyToQueries(p)
	if err != nil {
		return nil, errors.Wrap(err, "converting policy to query")
	}

	predicateTuples := make([]predicateTuple, 0, len(queries))
	for _, q := range queries {
		// Generate the deployment and image predicate
		imgPredicate, err := imageFactory.GeneratePredicate(q)
		if err != nil {
			return nil, err
		}
		deploymentPredicate, err := deploymentFactory.GeneratePredicate(q)
		if err != nil {
			return nil, err
		}

		processPredicate, err := processFactory.GeneratePredicate(q)
		if err != nil {
			return nil, err
		}
		predicateTuples = append(predicateTuples, predicateTuple{imagePredicate: imgPredicate, deploymentPredicate: deploymentPredicate, processPredicate: processPredicate})
	}

	return &matcherImpl{
		predicates: predicateTuples,
	}, nil
}
