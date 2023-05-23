package defaults

import (
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stretchr/testify/assert"
)

func TestIsDefaultRole(t *testing.T) {
	defaultRoleWithTraits := &storage.Role{Name: Admin, Traits: &storage.Traits{Origin: storage.Traits_DEFAULT}}
	defaultRoleWithoutTraits := &storage.Role{Name: Admin}
	nonDefaultRole := &storage.Role{Name: "some-random-role"}

	assert.True(t, IsDefaultRole(defaultRoleWithTraits))
	assert.True(t, IsDefaultRole(defaultRoleWithoutTraits))
	assert.False(t, IsDefaultRole(nonDefaultRole))
}

func TestIsDefaultAccessScope(t *testing.T) {
	defaultAccessScopeWithTraits := &storage.SimpleAccessScope{Id: AccessScopeIncludeAll.GetId(),
		Traits: &storage.Traits{Origin: storage.Traits_DEFAULT}}
	defaultAccessScopeWithoutTraits := &storage.SimpleAccessScope{Id: AccessScopeIncludeAll.GetId()}
	nonDefaultAccessScope := &storage.SimpleAccessScope{Id: "some-random-access-scope"}

	assert.True(t, IsDefaultAccessScope(defaultAccessScopeWithTraits))
	assert.True(t, IsDefaultAccessScope(defaultAccessScopeWithoutTraits))
	assert.False(t, IsDefaultAccessScope(nonDefaultAccessScope))
}

func TestRoleToPermSetMapping(t *testing.T) {
	defaultPermSets := GetDefaultPermissionSets()
	var permSetNames set.Set[string]
	for _, defaultPermSet := range defaultPermSets {
		permSetNames.Add(defaultPermSet.GetName())
	}

	for _, roleName := range DefaultRoleNames.AsSlice() {
		assert.True(t, permSetNames.Contains(roleName))
	}
}
