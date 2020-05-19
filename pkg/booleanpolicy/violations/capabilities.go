package violations

import (
	"errors"
	"fmt"
	"strings"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/booleanpolicy/augmentedobjs"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/searchbasedpolicies"
)

var (
	dropCapabilityFields = newViolationFieldSet(search.DropCapabilities.String(), augmentedobjs.ContainerNameCustomTag)
	addCapabilityFields  = newViolationFieldSet(search.AddCapabilities.String(), augmentedobjs.ContainerNameCustomTag)
)

func dropCapabilityPrinter(sectionName string, fieldMap map[string][]string) (*searchbasedpolicies.Violations, error) {
	if ok, err := dropCapabilityFields.validateMatchFields(fieldMap); !ok {
		return nil, err
	}
	if lenContainers := len(fieldMap[augmentedobjs.ContainerNameCustomTag]); lenContainers != 1 {
		return nil, fmt.Errorf("unexpected number of container names: %d", lenContainers)
	}
	var sb strings.Builder
	sb.WriteString("Container ")
	if containerName := fieldMap[augmentedobjs.ContainerNameCustomTag][0]; containerName != "" {
		fmt.Fprintf(&sb, "%s ", containerName)
	}
	sb.WriteString("adds ")
	switch capLen := len(fieldMap[search.DropCapabilities.String()]); {
	case capLen == 1:
		sb.WriteString("capability ")
	case capLen > 1:
		sb.WriteString("capabilities ")
	default:
		return nil, errors.New("Missing capabilities")
	}
	sb.WriteString(stringSliceToSortedSentence(fieldMap[search.DropCapabilities.String()]))
	return &searchbasedpolicies.Violations{AlertViolations: []*storage.Alert_Violation{{Message: sb.String()}}}, nil
}

func addCapabilityPrinter(sectionName string, fieldMap map[string][]string) (*searchbasedpolicies.Violations, error) {
	if ok, err := addCapabilityFields.validateMatchFields(fieldMap); !ok {
		return nil, err
	}
	if lenContainers := len(fieldMap[augmentedobjs.ContainerNameCustomTag]); lenContainers != 1 {
		return nil, fmt.Errorf("unexpected number of container names: %d", lenContainers)
	}
	var sb strings.Builder
	sb.WriteString("Container ")
	if containerName := fieldMap[augmentedobjs.ContainerNameCustomTag][0]; containerName != "" {
		fmt.Fprintf(&sb, "%s ", containerName)
	}
	sb.WriteString("adds ")
	switch capLen := len(fieldMap[search.AddCapabilities.String()]); {
	case capLen == 1:
		sb.WriteString("capability ")
	case capLen > 1:
		sb.WriteString("capabilities ")
	default:
		return nil, errors.New("Missing capabilities")
	}
	sb.WriteString(stringSliceToSortedSentence(fieldMap[search.AddCapabilities.String()]))
	return &searchbasedpolicies.Violations{AlertViolations: []*storage.Alert_Violation{{Message: sb.String()}}}, nil
}
