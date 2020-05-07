package booleanpolicy

import (
	"context"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/pathutil"
	"github.com/stackrox/rox/pkg/searchbasedpolicies"
)

const (
	imageAugmentKey   = "Image"
	processAugmentKey = "ProcessIndicator"
)

var (
	deploymentEvalFactory = evaluator.MustCreateNewFactory(
		pathutil.NewAugmentedObjMeta((*storage.Deployment)(nil)).
			AddPlainObjectAt([]string{"Containers", imageAugmentKey}, (*storage.Image)(nil)).
			AddAugmentedObjectAt(
				[]string{"Containers", processAugmentKey},
				pathutil.NewAugmentedObjMeta((*storage.ProcessIndicator)(nil)).
					AddPlainObjectAt([]string{"WhitelistStatus"}, (*whitelistResult)(nil)),
			),
	)

	imageEvalFactory = evaluator.MustCreateNewFactory(pathutil.NewAugmentedObjMeta((*storage.Image)(nil)))
)

type whitelistResult struct {
	NotWhitelisted bool `search:"Not Whitelisted"`
}

// An ImageMatcher matches images against a policy.
type ImageMatcher interface {
	MatchImage(ctx context.Context, image *storage.Image) (searchbasedpolicies.Violations, error)
}

// A DeploymentMatcher matches deployments against a policy.
type DeploymentMatcher interface {
	MatchDeployment(ctx context.Context, deployment *storage.Deployment, images []*storage.Image, pi *storage.ProcessIndicator) (searchbasedpolicies.Violations, error)
}

type sectionAndEvaluator struct {
	sectionName string
	evaluator   evaluator.Evaluator
}

// BuildDeploymentMatcher builds a matcher for deployments against the given policy,
// which must be a boolean policy.
func BuildDeploymentMatcher(p *storage.Policy) (DeploymentMatcher, error) {
	sectionsAndEvals, err := getSectionsAndEvals(&deploymentEvalFactory, p)
	if err != nil {
		return nil, err
	}

	return &matcherImpl{
		evaluators: sectionsAndEvals,
	}, nil
}

// BuildImageMatcher builds a matcher for images against the given policy,
// which must be a boolean policy.
func BuildImageMatcher(p *storage.Policy) (ImageMatcher, error) {
	sectionsAndEvals, err := getSectionsAndEvals(&imageEvalFactory, p)
	if err != nil {
		return nil, err
	}
	return &matcherImpl{evaluators: sectionsAndEvals}, nil
}

func getSectionsAndEvals(factory *evaluator.Factory, p *storage.Policy) ([]sectionAndEvaluator, error) {
	if err := Validate(p); err != nil {
		return nil, err
	}

	sectionsAndEvals := make([]sectionAndEvaluator, 0, len(p.GetPolicySections()))
	for _, section := range p.GetPolicySections() {
		sectionQ, err := sectionToQuery(section)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid section %q", section.GetSectionName())
		}
		eval, err := factory.GenerateEvaluator(sectionQ)
		if err != nil {
			return nil, errors.Wrapf(err, "generating evaluator for section %q", section.GetSectionName())
		}
		sectionsAndEvals = append(sectionsAndEvals, sectionAndEvaluator{section.GetSectionName(), eval})
	}

	return sectionsAndEvals, nil
}
