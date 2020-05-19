package violations

import (
	"fmt"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/searchbasedpolicies"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/utils"
)

// A ViolationPrinterFunc prints violation messages from a section name and map of fields to values
type ViolationPrinterFunc func(string, map[string][]string) (*searchbasedpolicies.Violations, error)

type violationFieldSet set.StringSet

func (v violationFieldSet) equal(other violationFieldSet) bool {
	return (set.StringSet)(v).Equal((set.StringSet)(other))
}

func newViolationFieldSet(fields ...string) *violationFieldSet {
	violationFields := set.NewStringSet(fields...)
	return (*violationFieldSet)(&violationFields)
}

func newViolationFieldSetFromMatchFields(fieldMap map[string][]string) violationFieldSet {
	keys := make([]string, 0, len(fieldMap))
	for key := range fieldMap {
		keys = append(keys, key)
	}
	return (violationFieldSet)(set.NewStringSet(keys...))
}

func (v violationFieldSet) validateMatchFields(fieldMap map[string][]string) (bool, error) {
	matchFields := newViolationFieldSetFromMatchFields(fieldMap)
	if !v.equal(matchFields) {
		return false, utils.Should(fmt.Errorf("mismatched fields printer: %v != match: %v", v, matchFields))
	}
	return true, nil
}

type violationPrinter struct {
	fields  *violationFieldSet
	printer ViolationPrinterFunc
}

func defaultViolationPrinter(sectionName string, fieldMap map[string][]string) (*searchbasedpolicies.Violations, error) {
	return &searchbasedpolicies.Violations{AlertViolations: []*storage.Alert_Violation{{Message: fmt.Sprintf("%+v", fieldMap)}}}, nil
}

var (
	// TODO(rc) make violationPrinter hashable for quick lookup based on fields
	printersAndFieldsByStage = map[storage.LifecycleStage][]*violationPrinter{
		storage.LifecycleStage_DEPLOY: {
			{dropCapabilityFields, dropCapabilityPrinter},
			{addCapabilityFields, addCapabilityPrinter},
		},
		storage.LifecycleStage_BUILD: {},
	}
)

func lookupViolationPrinter(stage storage.LifecycleStage, fieldMap map[string][]string) ViolationPrinterFunc {
	matchFields := newViolationFieldSetFromMatchFields(fieldMap)
	if printersAndFields, ok := printersAndFieldsByStage[stage]; ok {
		for _, p := range printersAndFields {
			if p.fields.equal(matchFields) {
				return p.printer
			}
		}
	}
	return defaultViolationPrinter
}

// ViolationPrinter creates violation messages based on evaluation results
func ViolationPrinter(stage storage.LifecycleStage, sectionName string, result *evaluator.Result) ([]*storage.Alert_Violation, error) {
	errorList := errorhelpers.NewErrorList("violation printer")
	violations := make([]*searchbasedpolicies.Violations, 0)
	for _, fieldMap := range result.Matches {
		printer := lookupViolationPrinter(stage, fieldMap)
		violation, err := printer(sectionName, fieldMap)
		if err != nil {
			errorList.AddError(err)
			continue
		}
		violations = append(violations, violation)
	}
	alertViolations := make([]*storage.Alert_Violation, 0)
	for _, violation := range violations {
		alertViolations = append(alertViolations, violation.AlertViolations...)
	}
	return alertViolations, errorList.ToError()
}
