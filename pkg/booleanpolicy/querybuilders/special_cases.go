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
		mappedValues := make([]string, 0, len(group.GetValues()))
		for _, value := range group.GetValues() {
			mappedValues = append(mappedValues, fmt.Sprintf("%s%s", basematchers.GreaterThanOrEqualTo, value.GetValue()))
		}
		return []*query.FieldQuery{{
			Field:  search.ServiceAccountPermissionLevel.String(),
			Values: mappedValues,
			Negate: group.GetNegate(),
		}}
	})
}
