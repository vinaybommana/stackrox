package manager

import (
	"strings"

	licenseproto "github.com/stackrox/rox/generated/shared/license"
)

func isStackRoxLicense(licenseMD *licenseproto.License_Metadata) bool {
	return strings.HasSuffix(licenseMD.GetLicensedForId(), "@stackrox.com")
}
