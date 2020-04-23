package booleanpolicy

import (
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/errorhelpers"
)

// Validate validates the policy, to make sure it's a well-formed Boolean policy.
func Validate(p *storage.Policy) error {
	errorList := errorhelpers.NewErrorList("policy validation")
	if p.GetPolicyVersion() != Version {
		errorList.AddStringf("invalid version for boolean policy (got %q)", p.GetPolicyVersion())
	}
	if p.GetName() == "" {
		errorList.AddString("no name specified")
	}
	return errorList.ToError()
}
