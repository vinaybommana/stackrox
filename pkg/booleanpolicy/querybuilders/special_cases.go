package querybuilders

import (
	"fmt"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/booleanpolicy/query"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/predicate/basematchers"
)

// ForK8sRBAC returns a specific query builder for K8s RBAC.
// Note that for K8s RBAC, the semantics are that
// the user specifies a value, and the policy matches if the actual permission
// is greater than or equal to that value.
func ForK8sRBAC() QueryBuilder {
	return queryBuilderFunc(func(group *storage.PolicyGroup) []*query.FieldQuery {
		return []*query.FieldQuery{{
			Field: search.ServiceAccountPermissionLevel.String(),
			Values: mapValues(group, func(s string) string {
				return fmt.Sprintf("%s%s", basematchers.GreaterThanOrEqualTo, s)
			}),
			Negate: group.GetNegate(),
		}}
	})
}

// ForDropCaps returns a specific query builder for drop capabilities.
// Note that here, we always negate -- the user specifies a list of capabilities that _must_ be dropped,
// so we want to find deployments that don't drop these capabilities.
func ForDropCaps() QueryBuilder {
	return queryBuilderFunc(func(group *storage.PolicyGroup) []*query.FieldQuery {
		return []*query.FieldQuery{{
			Field:    search.DropCapabilities.String(),
			Negate:   true,
			Values:   mapValues(group, valueToStringExact),
			Operator: operatorProtoMap[group.GetBooleanOperator()],
		}}
	})
}
