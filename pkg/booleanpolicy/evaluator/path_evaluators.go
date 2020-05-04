package evaluator

import (
	"fmt"
	"reflect"

	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/traverseutil"
	"github.com/stackrox/rox/pkg/search/fieldmap"
)

func wrapEvaluatorInPath(parentType reflect.Type, fieldPath fieldmap.FieldPath, evaluator internalEvaluator) (internalEvaluator, error) {
	if len(fieldPath) == 0 {
		return evaluator, nil
	}

	first, rest := fieldPath[0], fieldPath[1:]
	childEvaluator, err := wrapEvaluatorInPath(first.Type, rest, evaluator)
	if err != nil {
		return nil, err
	}

	// Wrap the predicate in field access.
	return traverseParentToField(parentType, first, childEvaluator)
}

func traverseParentToField(parentType reflect.Type, field reflect.StructField, evaluator internalEvaluator) (constructedEvaluator internalEvaluator, err error) {
	switch parentType.Kind() {
	case reflect.Array, reflect.Slice:
		return traverseSliceToField(parentType, field, evaluator)
	case reflect.Ptr:
		return traversePtrToField(parentType, field, evaluator)
	case reflect.Struct:
		return traverseStructToField(field, evaluator)
	default:
		return alwaysFalse, fmt.Errorf("cannot follow: %+v", field)
	}
}

func traverseSliceToField(parentType reflect.Type, field reflect.StructField, evaluator internalEvaluator) (internalEvaluator, error) {
	nestedEvaluator, err := traverseParentToField(parentType.Elem(), field, evaluator)
	if err != nil {
		return nil, err
	}

	return internalEvaluatorFunc(func(path *traverseutil.Path, instance reflect.Value) (*Result, bool) {
		length := instance.Len()
		if length == 0 {
			return nil, false
		}

		var results []*Result
		for i := 0; i < length; i++ {
			if res, matches := nestedEvaluator.Evaluate(path.WithSliceIndexed(i), instance.Index(i)); matches {
				results = append(results, res)
			}
		}
		if len(results) > 0 {
			return mergeResults(results), true
		}
		return nil, false
	}), nil

}

func traversePtrToField(parentType reflect.Type, field reflect.StructField, evaluator internalEvaluator) (internalEvaluator, error) {
	nestedEvaluator, err := traverseParentToField(parentType.Elem(), field, evaluator)
	if err != nil {
		return nil, err
	}

	return internalEvaluatorFunc(func(path *traverseutil.Path, instance reflect.Value) (*Result, bool) {
		if instance.IsNil() {
			return nil, false
		}
		return nestedEvaluator.Evaluate(path, instance.Elem())
	}), nil
}

func traverseStructToField(field reflect.StructField, evaluator internalEvaluator) (internalEvaluator, error) {
	return internalEvaluatorFunc(func(path *traverseutil.Path, instance reflect.Value) (*Result, bool) {
		nextValue := instance.FieldByIndex(field.Index)
		// We cannot traverse into a nil interface.
		if nextValue.Kind() == reflect.Interface && nextValue.IsNil() {
			return nil, false
		}
		return evaluator.Evaluate(path.WithFieldTraversed(field.Name), nextValue)
	}), nil
}
