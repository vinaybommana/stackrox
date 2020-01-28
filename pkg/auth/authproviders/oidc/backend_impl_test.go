package oidc

import "github.com/stackrox/rox/pkg/auth/authproviders"

var (
	_ authproviders.RefreshTokenEnabledBackend = (*backendImpl)(nil)
)
