package querybuilders

import (
	"fmt"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/booleanpolicy/augmentedobjs"
	"github.com/stackrox/rox/pkg/booleanpolicy/query"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/stringutils"
)

// ForCompound returns a custom query builder for a compound field that contains two values.
func ForCompound(field string) QueryBuilder {
	return queryBuilderFunc(func(group *storage.PolicyGroup) []*query.FieldQuery {
		return []*query.FieldQuery{
			{
				Field:    field,
				Operator: operatorProtoMap[group.GetBooleanOperator()],
				Values: mapValues(group, func(s string) string {
					first, second := stringutils.Split2(s, "=")
					// Compound fields are augmented and stored as "firstValue\tsecondValue"
					// To match this, we create the regex "(firstRegex)\t(secondRegex)",
					// replacing empty component by a ".*"
					return fmt.Sprintf("%s(%s)%s(%s)",
						search.RegexPrefix,
						stringutils.OrDefault(first, ".*"),
						augmentedobjs.CompositeFieldCharSep,
						stringutils.OrDefault(second, ".*"))
				}),
			},
		}
	})
}
