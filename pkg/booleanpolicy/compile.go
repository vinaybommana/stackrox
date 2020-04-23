package booleanpolicy

import (
	"github.com/pkg/errors"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
)

func policyToQueries(p *storage.Policy) ([]*v1.Query, error) {
	sectionQs := make([]*v1.Query, 0, len(p.GetPolicySections()))
	for _, section := range p.GetPolicySections() {
		sectionQ, err := sectionToQuery(section)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid section %q", section.GetSectionName())
		}
		sectionQs = append(sectionQs, sectionQ)
	}
	return sectionQs, nil
}

func sectionToQuery(section *storage.PolicySection) (*v1.Query, error) {
	groupQs := make([]*v1.Query, 0, len(section.GetPolicyGroups()))
	for _, group := range section.GetPolicyGroups() {
		groupQ, err := policyGroupToQuery(group)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid group for field %q", group.GetFieldName())
		}
		groupQs = append(groupQs, groupQ)
	}
	// TODO(viswa): handle linked fields.
	return search.ConjunctionQuery(groupQs...), nil
}

func policyGroupToQuery(group *storage.PolicyGroup) (*v1.Query, error) {
	if len(group.GetValues()) == 0 {
		return nil, errors.New("no values")
	}

	metadata := fieldsToQB[group.GetFieldName()]
	if metadata == nil {
		return nil, errors.New("no QB known for this group")
	}

	if metadata.negationForbidden && group.GetNegate() {
		return nil, errors.Errorf("invalid group: negation not allowed for field %s", group.GetFieldName())
	}
	if metadata.operatorsForbidden && len(group.GetValues()) != 1 {
		return nil, errors.Errorf("invalid group: operators not allowed for field %s", group.GetFieldName())
	}

	return metadata.qb.QueryForGroup(group), nil
}
