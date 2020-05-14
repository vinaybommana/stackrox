package booleanpolicy

import (
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/booleanpolicy/query"
)

func sectionToQuery(section *storage.PolicySection) (*query.Query, error) {
	fieldQueries := make([]*query.FieldQuery, 0, len(section.GetPolicyGroups()))
	for _, group := range section.GetPolicyGroups() {
		fqs, err := policyGroupToFieldQueries(group)
		if err != nil {
			return nil, err
		}
		fieldQueries = append(fieldQueries, fqs...)
	}
	return &query.Query{FieldQueries: fieldQueries}, nil
}

func policyGroupToFieldQueries(group *storage.PolicyGroup) ([]*query.FieldQuery, error) {
	if len(group.GetValues()) == 0 {
		return nil, errors.New("no values")
	}

	metadata := fieldsToQB[group.GetFieldName()]
	if metadata == nil || metadata.qb == nil {
		return nil, errors.Errorf("no QB known for group %q", group.GetFieldName())
	}

	if metadata.negationForbidden && group.GetNegate() {
		return nil, errors.Errorf("invalid group: negation not allowed for field %s", group.GetFieldName())
	}
	if metadata.operatorsForbidden && len(group.GetValues()) != 1 {
		return nil, errors.Errorf("invalid group: operators not allowed for field %s", group.GetFieldName())
	}

	fqs := metadata.qb.FieldQueriesForGroup(group)
	if len(fqs) == 0 {
		return nil, errors.New("invalid group: no queries formed")
	}

	return fqs, nil
}
