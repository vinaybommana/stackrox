package authn

import (
	"context"
	"errors"
	"time"

	"github.com/stackrox/rox/pkg/auth"
	"github.com/stackrox/rox/pkg/contextutil"
	"github.com/stackrox/rox/pkg/errox"
	"github.com/stackrox/rox/pkg/grpc/requestinfo"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/sync"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	// TODO: ROX-17312: Review and adjust logger parameters
	cacheSize          = 500
	rateLimitFrequency = 5 * time.Minute
	logBurstSize       = 5
)

var (
	once sync.Once
	log  *logging.RateLimitedLogger
)

type contextUpdater struct {
	extractor IdentityExtractor
}

func getRateLimitedLogger() *logging.RateLimitedLogger {
	once.Do(func() {
		// TODO: ROX-17312: Review and adjust logger parameters
		log = logging.NewRateLimitLogger(
			logging.LoggerForModule(),
			cacheSize,
			1,
			rateLimitFrequency,
			logBurstSize,
		)
	})
	return log
}

func (u contextUpdater) updateContext(ctx context.Context) (context.Context, error) {
	ri := requestinfo.FromContext(ctx)
	id, err := u.extractor.IdentityForRequest(ctx, ri)
	if err != nil {
		if errors.Is(err, jwt.ErrExpired) {
			getRateLimitedLogger().Debugf("Cannot extract identity: token expired")
		} else {
			getRateLimitedLogger().WarnL(ri.Hostname, "Cannot extract identity: %v", err)
		}
		// Ignore id value if error is not nil.
		return context.WithValue(ctx, identityErrorContextKey{}, errox.NoCredentials.CausedBy(err)), nil
	}
	if id != nil {
		// Only service identities can have no roles assigned.
		if len(id.Roles()) == 0 && id.Service() == nil {
			return context.WithValue(ctx, identityErrorContextKey{}, auth.ErrNoValidRole), nil
		}
		return context.WithValue(ctx, identityContextKey{}, id), nil
	}
	return ctx, nil
}

// ContextUpdater returns a context updater for the given identity extractors
func ContextUpdater(extractors ...IdentityExtractor) contextutil.ContextUpdater {
	return contextUpdater{extractor: CombineExtractors(extractors...)}.updateContext
}
