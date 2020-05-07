package evaluator

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/mapeval"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/pathutil"
	"github.com/stackrox/rox/pkg/booleanpolicy/query"
	"github.com/stackrox/rox/pkg/protoreflect"
	"github.com/stackrox/rox/pkg/readable"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/predicate/basematchers"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	timestampPtrType = reflect.TypeOf((*types.Timestamp)(nil))
)

// A baseEvaluator is an evaluator that operates on an individual field at the leaf of an object.
type baseEvaluator interface {
	Evaluate(*pathutil.Path, reflect.Value) (*Result, bool)
}

type baseEvaluatorFunc func(*pathutil.Path, reflect.Value) (*Result, bool)

func (f baseEvaluatorFunc) Evaluate(path *pathutil.Path, value reflect.Value) (*Result, bool) {
	return f(path, value)
}

func createBaseEvaluator(fieldName string, fieldType reflect.Type, values []string, negate bool, operator query.Operator) (baseEvaluator, error) {
	if len(values) == 0 {
		return nil, errors.New("no values in query")
	}
	if len(values) > 1 {
		if operator != query.Or && operator != query.And {
			return nil, errors.Errorf("invalid operator: %s", operator)
		}
	}
	kind := fieldType.Kind()
	generatorForKind, err := getMatcherGeneratorForKind(kind)
	if err != nil {
		return nil, err
	}

	baseMatchers := make([]baseMatcherAndExtractor, 0, len(values))
	for _, value := range values {
		m, err := generatorForKind(value, fieldType)
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

func combineMatchersIntoEvaluator(fieldName string, matchers []baseMatcherAndExtractor, operator query.Operator) baseEvaluator {
	return baseEvaluatorFunc(func(path *pathutil.Path, instance reflect.Value) (*Result, bool) {
		matchingValues := set.NewStringSet()
		var matches []string
		for _, m := range matchers {
			valuesAndMatches := m(instance)
			// This means there were no values.
			if len(valuesAndMatches) == 0 {
				return nil, false
			}
			var atLeastOneSuccess bool
			for _, valueAndMatch := range valuesAndMatches {
				if valueAndMatch.matched {
					if matchingValues.Add(valueAndMatch.value) {
						matches = append(matches, valueAndMatch.value)
					}
					atLeastOneSuccess = true
				}
			}
			// If not matched, and it's an And, then we can early exit.
			if !atLeastOneSuccess && operator == query.And {
				return nil, false
			}
		}
		if matchingValues.Cardinality() == 0 {
			return nil, false
		}
		return resultWithSingleMatch(fieldName, path, matches...), true
	})
}

func combineMatchersIntoEvaluatorNegated(fieldName string, matchers []baseMatcherAndExtractor, operator query.Operator) baseEvaluator {
	return baseEvaluatorFunc(func(path *pathutil.Path, instance reflect.Value) (*Result, bool) {
		matchingValues := set.NewStringSet()
		for _, m := range matchers {
			valuesAndMatches := m(instance)
			// This means there were no values.
			if len(valuesAndMatches) == 0 {
				return nil, false
			}
			var atLeastOneSuccess bool
			for _, valueAndMatch := range valuesAndMatches {
				if !valueAndMatch.matched {
					matchingValues.Add(valueAndMatch.value)
					atLeastOneSuccess = true
				}
			}

			// If not matched, and it's an Or, then we can early exit.
			// Since we're negating, this check is correct by de Morgan's law.
			// !(A OR B) <=> !A AND !B, therefore if operator is OR and A _does_ match,
			// we can conclude that !A is false => !A AND !B is false => !(A OR B) is false.
			if !atLeastOneSuccess && operator == query.Or {
				return nil, false
			}
		}
		if matchingValues.Cardinality() == 0 {
			return nil, false
		}
		return resultWithSingleMatch(fieldName, path, matchingValues.AsSlice()...), true
	})
}

func getMatcherGeneratorForKind(kind reflect.Kind) (baseMatcherGenerator, error) {
	switch kind {
	case reflect.String:
		return generateStringMatcher, nil
	case reflect.Ptr:
		return generatePtrMatcher, nil
	case reflect.Array, reflect.Slice:
		return generateSliceMatcher, nil
	case reflect.Map:
		return generateMapMatcher, nil
	case reflect.Bool:
		return generateBoolMatcher, nil
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		return generateIntMatcher, nil
	case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint:
		return generateUintMatcher, nil
	case reflect.Float64, reflect.Float32:
		return generateFloatMatcher, nil
	default:
		return nil, errors.Errorf("invalid kind for base query: %s", kind)
	}
}

type baseMatcherGenerator func(string, reflect.Type) (baseMatcherAndExtractor, error)

// A baseMatcherAndExtractor takes a value of a given type, extracts a human-readable string value
// and returns whether it matched or not.
// IMPORTANT: in every valueMatchedPair, value _must_ be returned even if _matched_ is false,
// since that enables us to get the value out even if the caller is going to negate this query.
type baseMatcherAndExtractor func(reflect.Value) []valueMatchedPair

type valueMatchedPair struct {
	value   string
	matched bool
}

func resultWithSingleMatch(fieldName string, path *pathutil.Path, values ...string) *Result {
	return &Result{Matches: map[string][]Match{fieldName: {{Path: path, Values: values}}}}
}

func generateStringMatcher(value string, _ reflect.Type) (baseMatcherAndExtractor, error) {
	baseMatcher, err := basematchers.ForString(value)
	if err != nil {
		return nil, err
	}
	return func(instance reflect.Value) []valueMatchedPair {
		if instance.Kind() != reflect.String {
			return nil
		}
		asStr := instance.String()
		return []valueMatchedPair{{value: asStr, matched: baseMatcher(asStr)}}
	}, nil
}

func generateSliceMatcher(value string, fieldType reflect.Type) (baseMatcherAndExtractor, error) {
	underlyingType := fieldType.Elem()
	matcherGenerator, err := getMatcherGeneratorForKind(underlyingType.Kind())
	if err != nil {
		return nil, err
	}
	subMatcher, err := matcherGenerator(value, underlyingType)
	if err != nil {
		return nil, err
	}

	return func(instance reflect.Value) []valueMatchedPair {
		length := instance.Len()
		if length == 0 {
			// An empty slice matches no queries, but we want to bubble this up,
			// for callers that are negating.
			return []valueMatchedPair{{value: "<empty>", matched: false}}
		}
		valuesAndMatches := make([]valueMatchedPair, 0, length)
		for i := 0; i < length; i++ {
			valuesAndMatches = append(valuesAndMatches, subMatcher(instance.Index(i))...)
		}
		return valuesAndMatches
	}, nil
}

func generateTimestampMatcher(value string) (baseMatcherAndExtractor, error) {
	var baseMatcher func(*types.Timestamp) bool
	if value != search.NullString {
		var err error
		baseMatcher, err = basematchers.ForTimestamp(value)
		if err != nil {
			return nil, err
		}
	}
	return func(instance reflect.Value) []valueMatchedPair {
		ts, ok := instance.Interface().(*types.Timestamp)
		if !ok {
			return nil
		}
		if ts == nil {
			if value == search.NullString {
				return []valueMatchedPair{{value: "<empty timestamp>", matched: true}}
			}
			return nil
		}
		return []valueMatchedPair{{value: readable.ProtoTime(ts), matched: value != "-" && baseMatcher(ts)}}
	}, nil
}

func generatePtrMatcher(value string, fieldType reflect.Type) (baseMatcherAndExtractor, error) {
	// Special case for pointer to timestamp.
	if fieldType == timestampPtrType {
		return generateTimestampMatcher(value)
	}

	underlyingType := fieldType.Elem()
	matcherGenerator, err := getMatcherGeneratorForKind(underlyingType.Kind())
	if err != nil {
		return nil, err
	}
	subMatcher, err := matcherGenerator(value, underlyingType)
	if err != nil {
		return nil, err
	}

	return func(instance reflect.Value) []valueMatchedPair {
		if instance.IsNil() {
			return []valueMatchedPair{{value: "<nil>", matched: value == search.NullString}}
		}
		subMatches := subMatcher(instance.Elem())
		// If the value is null, and the pointer is not nil, it did not match.
		// So just use the values from the subMatcher but always set matched
		// to false.
		if value == search.NullString {
			for i := range subMatches {
				subMatches[i].matched = false
			}
		}
		return subMatches
	}, nil
}

func generateBoolMatcher(value string, _ reflect.Type) (baseMatcherAndExtractor, error) {
	baseMatcher, err := basematchers.ForBool(value)
	if err != nil {
		return nil, err
	}
	return func(instance reflect.Value) []valueMatchedPair {
		if instance.Kind() != reflect.Bool {
			return nil
		}
		asBool := instance.Bool()
		return []valueMatchedPair{{value: fmt.Sprintf("%t", asBool), matched: baseMatcher(asBool)}}
	}, nil
}

func generateIntMatcher(value string, fieldType reflect.Type) (baseMatcherAndExtractor, error) {
	if enum, ok := reflect.Zero(fieldType).Interface().(protoreflect.ProtoEnum); ok {
		return generateEnumMatcher(value, enum)
	}

	baseMatcher, err := basematchers.ForInt(value)
	if err != nil {
		return nil, err
	}

	return func(instance reflect.Value) []valueMatchedPair {
		switch instance.Kind() {
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			asInt := instance.Int()
			return []valueMatchedPair{{value: fmt.Sprintf("%d", asInt), matched: baseMatcher(asInt)}}
		}
		return nil
	}, nil
}

func generateUintMatcher(value string, _ reflect.Type) (baseMatcherAndExtractor, error) {
	baseMatcher, err := basematchers.ForUint(value)
	if err != nil {
		return nil, err
	}

	return func(instance reflect.Value) []valueMatchedPair {
		switch instance.Kind() {
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			asUint := instance.Uint()
			return []valueMatchedPair{{value: fmt.Sprintf("%d", asUint), matched: baseMatcher(asUint)}}
		}
		return nil
	}, nil
}

func generateFloatMatcher(value string, _ reflect.Type) (baseMatcherAndExtractor, error) {
	baseMatcher, err := basematchers.ForFloat(value)
	if err != nil {
		return nil, err
	}

	return func(instance reflect.Value) []valueMatchedPair {
		switch instance.Kind() {
		case reflect.Float32, reflect.Float64:
			asFloat := instance.Float()
			return []valueMatchedPair{{value: fmt.Sprintf("%g", asFloat), matched: baseMatcher(asFloat)}}
		}
		return nil
	}, nil
}

func generateEnumMatcher(value string, enumRef protoreflect.ProtoEnum) (baseMatcherAndExtractor, error) {
	baseMatcher, numberToName, err := basematchers.ForEnum(value, enumRef)
	if err != nil {
		return nil, err
	}

	return func(instance reflect.Value) []valueMatchedPair {
		if instance.Kind() != reflect.Int32 {
			return nil
		}
		asInt := instance.Int()
		matchedValue := numberToName[int32(asInt)]
		if matchedValue == "" {
			utils.Should(errors.Errorf("enum query matched (%s), but no value in numberToName (%v) (got number: %d)",
				value, numberToName, asInt))
			matchedValue = strconv.Itoa(int(asInt))
		}
		return []valueMatchedPair{{value: matchedValue, matched: baseMatcher(asInt)}}
	}, nil
}

func generateMapMatcher(value string, _ reflect.Type) (baseMatcherAndExtractor, error) {
	baseMatcher, err := mapeval.Matcher(value)
	if err != nil {
		return nil, err
	}

	return func(instance reflect.Value) []valueMatchedPair {
		if instance.Kind() != reflect.Map {
			return nil
		}

		iter := instance.MapRange()
		return []valueMatchedPair{{value: "", matched: baseMatcher(iter)}}
	}, nil
}
