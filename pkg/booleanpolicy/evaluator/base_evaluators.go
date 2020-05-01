package evaluator

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/traverseutil"
	"github.com/stackrox/rox/pkg/booleanpolicy/query"
	"github.com/stackrox/rox/pkg/search/predicate/basematchers"
	"github.com/stackrox/rox/pkg/set"
)

func createBaseEvaluator(fieldName string, fieldType reflect.Type, values []string, negate bool, operator query.Operator) (internalEvaluator, error) {
	if len(values) == 0 {
		return nil, errors.New("no values in query")
	}
	if len(values) > 1 {
		if operator != query.Or && operator != query.And {
			return nil, errors.Errorf("invalid operator: %s", operator)
		}
	}
	kind := fieldType.Kind()
	generatorForKind := getMatcherGeneratorForKind(kind)
	if generatorForKind == nil {
		return nil, errors.Errorf("unknown field kind: %v", kind)
	}
	baseMatchers := make([]baseMatcherAndExtractor, 0, len(values))
	for _, value := range values {
		m, err := generatorForKind(value)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid value: %s for field %s", value, fieldName)
		}
		baseMatchers = append(baseMatchers, m)
	}

	if negate {
		return combineMatchersIntoEvaluatorNegated(fieldName, baseMatchers, operator), nil
	}
	return combineMatchersIntoEvaluator(fieldName, baseMatchers, operator), nil
}

func combineMatchersIntoEvaluator(fieldName string, matchers []baseMatcherAndExtractor, operator query.Operator) internalEvaluator {
	return internalEvaluatorFunc(func(path *traverseutil.Path, instance reflect.Value) (*Result, bool) {
		matchingValues := set.NewStringSet()
		for _, m := range matchers {
			value, matched := m(instance)
			if matched {
				matchingValues.Add(value)
				continue
			}
			// If not matched, and it's an And, then we can early exit.
			if operator == query.And {
				return nil, false
			}
		}
		if matchingValues.Cardinality() == 0 {
			return nil, false
		}
		return resultWithSingleMatch(fieldName, path, matchingValues.AsSlice()...), true
	})
}

func combineMatchersIntoEvaluatorNegated(fieldName string, matchers []baseMatcherAndExtractor, operator query.Operator) internalEvaluator {
	return internalEvaluatorFunc(func(path *traverseutil.Path, instance reflect.Value) (*Result, bool) {
		matchingValues := set.NewStringSet()
		for _, m := range matchers {
			value, matched := m(instance)
			if !matched {
				matchingValues.Add(value)
				continue
			}
			// If not matched, and it's an Or, then we can early exit.
			// Since we're negating, this check is correct by de Morgan's law.
			// !(A OR B) <=> !A AND !B, therefore if operator is OR and A _does_ match,
			// we can conclude that !A is false => !A AND !B is false => !(A OR B) is false.
			if operator == query.Or {
				return nil, false
			}
		}
		if matchingValues.Cardinality() == 0 {
			return nil, false
		}
		return resultWithSingleMatch(fieldName, path, matchingValues.AsSlice()...), true
	})
}

func getMatcherGeneratorForKind(kind reflect.Kind) baseMatcherGenerator {
	switch kind {
	case reflect.String:
		return generateStringMatcher
	}
	return nil
}

type baseMatcherGenerator func(string) (baseMatcherAndExtractor, error)

// A baseMatcherAndExtractor takes a value of a given type, extracts a human-readable string value
// and returns whether it matched or not.
// IMPORTANT: it returns the string value irrespective of whether the value actually matched (ie, the second return value),
// since that enables us to get the value out even if the caller is going to negate this query.
type baseMatcherAndExtractor func(reflect.Value) (string, bool)

func resultWithSingleMatch(fieldName string, path *traverseutil.Path, values ...string) *Result {
	return &Result{Matches: map[string][]Match{fieldName: {{Path: path, Values: values}}}}
}

func generateStringMatcher(value string) (baseMatcherAndExtractor, error) {
	baseMatcher, err := basematchers.ForString(value)
	if err != nil {
		return nil, err
	}
	return func(instance reflect.Value) (string, bool) {
		if instance.Kind() != reflect.String {
			return "", false
		}
		asStr := instance.String()
		return asStr, baseMatcher(asStr)
	}, nil
}
