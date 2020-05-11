package querybuilders

import (
	"strings"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/booleanpolicy/query"
	"github.com/stackrox/rox/pkg/search"
)

var (
	operatorProtoMap = map[storage.BooleanOperator]query.Operator{
		storage.BooleanOperator_OR:  query.Or,
		storage.BooleanOperator_AND: query.And,
	}
)

func valueToStringExact(value string) string {
	return search.ExactMatchString(value)
}

func valueToStringRegex(value string) string {
	if strings.HasPrefix(value, search.RegexPrefix) {
		return value
	}
	return search.RegexPrefix + value
}

func mapValues(group *storage.PolicyGroup, f func(string) string) []string {
	out := make([]string, 0, len(group.GetValues()))
	for _, v := range group.GetValues() {
		var mappedValue string
		if f != nil {
			mappedValue = f(v.GetValue())
		} else {
			mappedValue = v.GetValue()
		}
		out = append(out, mappedValue)
	}
	return out
}

// A QueryBuilder builds queries for a specific policy group.
type QueryBuilder interface {
	FieldQueriesForGroup(group *storage.PolicyGroup) []*query.FieldQuery
}

type queryBuilderFunc func(group *storage.PolicyGroup) []*query.FieldQuery

func (f queryBuilderFunc) FieldQueriesForGroup(group *storage.PolicyGroup) []*query.FieldQuery {
	return f(group)
}

type fieldLabelQueryBuilder struct {
	fieldLabel   search.FieldLabel
	valueMapFunc func(string) string
}

func (f *fieldLabelQueryBuilder) FieldQueriesForGroup(group *storage.PolicyGroup) []*query.FieldQuery {
	fq := &query.FieldQuery{
		Field:    f.fieldLabel.String(),
		Values:   mapValues(group, f.valueMapFunc),
		Operator: operatorProtoMap[group.GetBooleanOperator()],
		Negate:   group.GetNegate(),
	}
	return []*query.FieldQuery{fq}
}

// ForFieldLabelExact returns a query builder that simply queries for the exact field value with the given search field label.
func ForFieldLabelExact(label search.FieldLabel) QueryBuilder {
	return &fieldLabelQueryBuilder{fieldLabel: label, valueMapFunc: valueToStringExact}
}

// ForFieldLabel returns a query builder that does a prefix match for the field value with the given search field label.
func ForFieldLabel(label search.FieldLabel) QueryBuilder {
	return &fieldLabelQueryBuilder{fieldLabel: label}
}

// ForFieldLabelRegex is like ForFieldLabel, but does a regex match.
func ForFieldLabelRegex(label search.FieldLabel) QueryBuilder {
	return &fieldLabelQueryBuilder{fieldLabel: label, valueMapFunc: valueToStringRegex}
}

// ForFieldLabelUpper is like ForFieldLabel, but does a match after converting the query to upper-case.
func ForFieldLabelUpper(label search.FieldLabel) QueryBuilder {
	return &fieldLabelQueryBuilder{fieldLabel: label, valueMapFunc: strings.ToUpper}
}
