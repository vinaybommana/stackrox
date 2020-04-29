package evaluator

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/traverseutil"
	"github.com/stackrox/rox/pkg/regexutils"
	"github.com/stackrox/rox/pkg/search"
)

func createBaseEvaluator(fieldType reflect.Type, value string) (internalEvaluator, error) {
	switch fieldType.Kind() {
	case reflect.String:
		return createStringEvaluator(value)
	}
	return nil, errors.New("unrecognized field")
}

func resultWithSingleMatch(path *traverseutil.Path, value interface{}) *Result {
	return &Result{Matches: []Match{{Path: path, Value: value}}}
}

func createStringEvaluator(value string) (internalEvaluator, error) {
	negated := strings.HasPrefix(value, search.NegationPrefix)
	if negated {
		value = strings.TrimPrefix(value, search.NegationPrefix)
	}
	if strings.HasPrefix(value, search.RegexPrefix) {
		value = strings.TrimPrefix(value, search.RegexPrefix)
		return stringRegexPredicate(value, negated)
	} else if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) && len(value) > 1 {
		return stringExactPredicate(value[1:len(value)-1], negated)
	}
	return stringPrefixPredicate(value, negated)
}

func stringRegexPredicate(value string, negated bool) (internalEvaluator, error) {
	matcher, err := regexp.Compile(value)
	if err != nil {
		return nil, err
	}
	return wrapStringEvaluator(func(path *traverseutil.Path, instance string) (*Result, bool) {
		// matched == negated is equivalent to !(matched XOR negated), which is what we want here
		if regexutils.MatchWholeString(matcher, instance) == negated {
			return nil, false
		}

		return resultWithSingleMatch(path, instance), true
	}), nil
}

func stringExactPredicate(value string, negated bool) (internalEvaluator, error) {
	return wrapStringEvaluator(func(path *traverseutil.Path, instance string) (*Result, bool) {
		// matched == negated is equivalent to !(matched XOR negated), which is what we want here
		if (instance == value) == negated {
			return nil, false
		}
		return resultWithSingleMatch(path, instance), true
	}), nil
}

func stringPrefixPredicate(value string, negated bool) (internalEvaluator, error) {
	return wrapStringEvaluator(func(path *traverseutil.Path, instance string) (*Result, bool) {
		// matched == negated is equivalent to !(matched XOR negated), which is what we want here
		if (value == search.WildcardString || strings.HasPrefix(instance, value)) == negated {
			return nil, false
		}
		return resultWithSingleMatch(path, instance), true
	}), nil
}

func wrapStringEvaluator(evaluator func(*traverseutil.Path, string) (*Result, bool)) internalEvaluator {
	return internalEvaluatorFunc(func(path *traverseutil.Path, instance reflect.Value) (*Result, bool) {
		if instance.Kind() != reflect.String {
			return nil, false
		}
		return evaluator(path, instance.String())
	})
}
