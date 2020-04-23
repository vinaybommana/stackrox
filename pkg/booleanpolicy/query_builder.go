package booleanpolicy

import (
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
)

type queryBuilder interface {
	// QueryForGroup generates a query based on a policy group, for a specific field name.
	// QueryForGroup implementations can assume that the group being passed has the expected field name for the query builder
	// and that the group has at least one value.
	// If callers do not respect this contract, the result is undefined, but will most likely be a runtime panic.
	QueryForGroup(group *storage.PolicyGroup) *v1.Query
}

type fieldLabelBasedQueryBuilder struct {
	fieldName  string
	fieldLabel search.FieldLabel
}

func (f *fieldLabelBasedQueryBuilder) QueryForGroup(group *storage.PolicyGroup) *v1.Query {
	matchFieldQs := make([]*v1.Query, 0, len(group.GetValues()))
	for _, value := range group.GetValues() {
		matchFieldQs = append(matchFieldQs, search.MatchFieldQuery(f.fieldLabel.String(), value.GetValue(), false))
	}
	var combinedQ *v1.Query
	if len(matchFieldQs) == 1 {
		combinedQ = matchFieldQs[0]
	} else if group.GetBooleanOperator() == storage.BooleanOperator_AND {
		combinedQ = search.ConjunctionQuery(matchFieldQs...)
	} else {
		combinedQ = search.DisjunctionQuery(matchFieldQs...)
	}
	if group.GetNegate() {
		// TODO(viswa): Figure out if the EmptyQuery() is necessary. It's not immediately clear based on the predicate logic.
		return search.NewBooleanQuery(&v1.ConjunctionQuery{Queries: []*v1.Query{search.EmptyQuery()}}, &v1.DisjunctionQuery{Queries: []*v1.Query{combinedQ}})
	}
	return combinedQ
}
