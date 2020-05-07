package evaluator

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/pathutil"
)

func wrapBaseEvaluator(baseEvaluator baseEvaluator) Evaluator {
	return evaluatorFunc(func(value pathutil.AugmentedValue) (*Result, bool) {
		return baseEvaluator.Evaluate(value.PathFromRoot(), value.Underlying())
	})
}

func (f *Factory) wrapBaseEvaluatorWithPathTraversal(pathToBase pathutil.MetaPath, baseEvaluator baseEvaluator) (Evaluator, error) {
	return wrapEvaluatorWithTraversal(f.rootType, pathToBase, wrapBaseEvaluator(baseEvaluator))
}

func wrapEvaluatorWithTraversal(currentType reflect.Type, pathToEvaluator pathutil.MetaPath, evaluator Evaluator) (Evaluator, error) {
	// Base case
	if len(pathToEvaluator) == 0 {
		return evaluator, nil
	}

	firstStep, remainingPath := pathToEvaluator[0], pathToEvaluator[1:]
	childEvaluator, err := wrapEvaluatorWithTraversal(firstStep.Type, remainingPath, evaluator)
	if err != nil {
		return nil, err
	}
	return takeMetaStep(currentType, firstStep, childEvaluator)
}

func takeMetaStep(currentType reflect.Type, metaStep pathutil.MetaStep, evaluator Evaluator) (Evaluator, error) {
	switch currentType.Kind() {
	case reflect.Array, reflect.Slice:
		return takeSliceMetaStep(currentType, metaStep, evaluator)
	case reflect.Ptr:
		return takePtrMetaStep(currentType, metaStep, evaluator)
	case reflect.Struct:
		return takeStructMetaStep(metaStep, evaluator), nil
	case reflect.Interface:
		return takeInterfaceMetaStep(metaStep, evaluator), nil
	default:
		return nil, errors.Errorf("cannot follow: %+v", metaStep)
	}
}

func takeInterfaceMetaStep(metaStep pathutil.MetaStep, evaluator Evaluator) Evaluator {
	return evaluatorFunc(func(value pathutil.AugmentedValue) (*Result, bool) {
		underlying := value.Underlying()
		if underlying.IsNil() {
			return nil, false
		}
		nextValue := value.Elem()
		if nextValue.Underlying().Kind() == reflect.Ptr {
			nextValue = nextValue.Elem()
		}
		if nextValue.Underlying().Kind() != reflect.Struct {
			return nil, false
		}
		nextValue, found := nextValue.TakeStep(metaStep)
		if !found {
			return nil, false
		}
		return evaluator.Evaluate(nextValue)
	})
}

func takeStructMetaStep(metaStep pathutil.MetaStep, evaluator Evaluator) Evaluator {
	return evaluatorFunc(func(value pathutil.AugmentedValue) (*Result, bool) {
		nextValue, found := value.TakeStep(metaStep)
		if !found {
			return nil, false
		}
		return evaluator.Evaluate(nextValue)
	})
}

func takeSliceMetaStep(currentType reflect.Type, metaStep pathutil.MetaStep, evaluator Evaluator) (Evaluator, error) {
	nestedEvaluator, err := takeMetaStep(currentType.Elem(), metaStep, evaluator)
	if err != nil {
		return nil, err
	}

	return evaluatorFunc(func(value pathutil.AugmentedValue) (*Result, bool) {
		length := value.Underlying().Len()
		if length == 0 {
			return nil, false
		}

		var results []*Result
		for i := 0; i < length; i++ {
			if res, matches := nestedEvaluator.Evaluate(value.Index(i)); matches {
				results = append(results, res)
			}
		}
		if len(results) > 0 {
			return mergeResults(results), true
		}
		return nil, false
	}), nil

}

func takePtrMetaStep(currentType reflect.Type, metaStep pathutil.MetaStep, evaluator Evaluator) (Evaluator, error) {
	nextStepEvaluator, err := takeMetaStep(currentType.Elem(), metaStep, evaluator)
	if err != nil {
		return nil, err
	}

	return evaluatorFunc(func(value pathutil.AugmentedValue) (*Result, bool) {
		if value.Underlying().IsNil() {
			return nil, false
		}
		return nextStepEvaluator.Evaluate(value.Elem())
	}), nil
}
