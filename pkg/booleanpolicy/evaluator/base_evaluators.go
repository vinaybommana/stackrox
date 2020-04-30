package evaluator

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/traverseutil"
	"github.com/stackrox/rox/pkg/search/predicate/basematchers"
)

func createBaseEvaluator(fieldName string, fieldType reflect.Type, value string) (internalEvaluator, error) {
	switch fieldType.Kind() {
	case reflect.String:
		return createStringEvaluator(fieldName, value)
	}
	return nil, errors.New("unrecognized field")
}

func resultWithSingleMatch(fieldName string, path *traverseutil.Path, value interface{}) *Result {
	return &Result{Matches: map[string][]Match{fieldName: {{Path: path, Value: value}}}}
}

func createStringEvaluator(fieldName, value string) (internalEvaluator, error) {
	baseMatcher, err := basematchers.ForString(value)
	if err != nil {
		return nil, err
	}
	return wrapStringMatcher(fieldName, baseMatcher), nil
}

func wrapStringMatcher(fieldName string, matcher func(string) bool) internalEvaluator {
	return internalEvaluatorFunc(func(path *traverseutil.Path, instance reflect.Value) (*Result, bool) {
		if instance.Kind() != reflect.String {
			return nil, false
		}
		instanceAsString := instance.String()
		if matcher(instanceAsString) {
			return resultWithSingleMatch(fieldName, path, instanceAsString), true
		}
		return nil, false
	})
}
