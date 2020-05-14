package booleanpolicy

import (
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/set"
)

var (
	runtimeFields  = set.NewFrozenStringSet(ProcessName, ProcessArguments, ProcessAncestor, ProcessUID, WhitelistsEnabled)
	whitelistField = set.NewFrozenStringSet(WhitelistsEnabled)
)

func policyHasFields(policy *storage.Policy, fieldSet set.FrozenStringSet) bool {
	for _, section := range policy.GetPolicySections() {
		for _, group := range section.GetPolicyGroups() {
			if fieldSet.Contains(group.GetFieldName()) {
				return true
			}
		}
	}
	return false
}

// ContainsRuntimeFields returns whether the policy contains runtime specific fields.
func ContainsRuntimeFields(policy *storage.Policy) bool {
	return policyHasFields(policy, runtimeFields)
}

// IsWhitelistEnabled returns whether a boolean policy has a policy group with the given name
func IsWhitelistEnabled(policy *storage.Policy) bool {
	return policyHasFields(policy, whitelistField)
}
