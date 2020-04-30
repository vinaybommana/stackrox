package booleanpolicy

import "github.com/stackrox/rox/generated/storage"

// IsWhitelistEnabled returns whether a boolean policy has a policy group with the given name
func IsWhitelistEnabled(policy *storage.Policy) bool {
	for _, section := range policy.GetPolicySections() {
		for _, group := range section.GetPolicyGroups() {
			if group.GetFieldName() == WhitelistsEnabled {
				return true
			}
		}
	}
	return false
}
