// +build !release

package singleton

import (
	"time"

	"github.com/stackrox/rox/pkg/license/publickeys"
	"github.com/stackrox/rox/pkg/license/validator"
	"github.com/stackrox/rox/pkg/timeutil"
)

func init() {

	registerValidatorRegistrationArgs(
		validatorRegistrationArgs{
			publickeys.Dev,
			func() validator.SigningKeyRestrictions {
				return validator.SigningKeyRestrictions{
					EarliestNotValidBefore:                  timeutil.MustParse(time.RFC3339, "2019-12-01T00:00:00Z"),
					LatestNotValidAfter:                     timeutil.MustParse(time.RFC3339, "2020-04-01T00:00:00Z"),
					MaxDuration:                             30 * 24 * time.Hour,
					AllowOffline:                            true,
					MaxNodeLimit:                            50,
					BuildFlavors:                            []string{"development"},
					AllowNoDeploymentEnvironmentRestriction: true,
				}
			},
		},
		// OLD VERSION - NO LONGER USED FOR NEW LICENSES
		validatorRegistrationArgs{
			publickeys.DevOld,
			func() validator.SigningKeyRestrictions {
				return validator.SigningKeyRestrictions{
					EarliestNotValidBefore:                  timeutil.MustParse(time.RFC3339, "2019-09-01T00:00:00Z"),
					LatestNotValidAfter:                     timeutil.MustParse(time.RFC3339, "2020-01-01T00:00:00Z"),
					MaxDuration:                             30 * 24 * time.Hour,
					AllowOffline:                            true,
					MaxNodeLimit:                            50,
					BuildFlavors:                            []string{"development"},
					AllowNoDeploymentEnvironmentRestriction: true,
				}
			},
		},
	)
}
