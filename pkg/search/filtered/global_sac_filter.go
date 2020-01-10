package filtered

import (
	"context"

	"github.com/stackrox/rox/pkg/sac"
)

type globalFilterImpl struct {
	scopeChecker sac.ScopeChecker
}

func (f *globalFilterImpl) Apply(ctx context.Context, from ...string) ([]string, error) {
	if ok, err := f.scopeChecker.Allowed(ctx); err != nil || !ok {
		return nil, err
	}
	return from, nil
}
