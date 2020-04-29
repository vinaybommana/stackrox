package evaluator

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/traverseutil"
	"github.com/stackrox/rox/pkg/booleanpolicy/query"
	"github.com/stackrox/rox/pkg/search/fieldmap"
	"github.com/stackrox/rox/pkg/utils"
)

// An Evaluator evalutes an object, and produces a result.
type Evaluator interface {
	Evaluate(instance interface{}) (*Result, bool)
}

// evaluatorFunc wraps a regular function as a predicate.
type evaluatorFunc func(instance interface{}) (*Result, bool)

func (f evaluatorFunc) Evaluate(instance interface{}) (*Result, bool) {
	return f(instance)
}

type internalEvaluator interface {
	Evaluate(*traverseutil.Path, reflect.Value) (*Result, bool)
}

type internalEvaluatorFunc func(*traverseutil.Path, reflect.Value) (*Result, bool)

func (f internalEvaluatorFunc) Evaluate(path *traverseutil.Path, value reflect.Value) (*Result, bool) {
	return f(path, value)
}

// Factory object stores the specs for each when walking the query.
type Factory struct {
	searchFields fieldmap.FieldMap
	rootType     reflect.Type
}

// NewFactory returns a new evaluator factory for the type of the given object.
func NewFactory(obj interface{}) Factory {
	return Factory{
		searchFields: fieldmap.MapSearchTagsToFieldPaths(obj),
		rootType:     reflect.TypeOf(obj),
	}
}

// GenerateEvaluator generates an evaluator that will evaluate the given query
// on objects of the factory's type.
func (f *Factory) GenerateEvaluator(q query.Query) (Evaluator, error) {
	internal, err := f.generateInternalEvaluator(q)
	if err != nil {
		return nil, err
	}
	return wrapInternal(internal), nil
}

func wrapInternal(ie internalEvaluator) Evaluator {
	return evaluatorFunc(func(in interface{}) (res *Result, matched bool) {
		defer func() {
			// Panics can occur in evaluators, mainly due to incorrect uses of reflect.
			// This is always a programming error, but let's not panic in prod over it.
			if r := recover(); r != nil {
				utils.Should(errors.Errorf("panic running evaluator: %v", r))
				res = nil
				matched = false
			}
		}()
		return ie.Evaluate(&traverseutil.Path{}, reflect.ValueOf(in))
	})
}

func (f *Factory) generateInternalEvaluator(q query.Query) (internalEvaluator, error) {
	switch underlying := q.GetUnderlying().(type) {
	case *query.Disjunction:
		return f.generateDisjunction(underlying)
	case *query.Conjunction:
		return f.generateConjunction(underlying)
	case *query.LinkedConjunction:
		return f.generateLinkedConjunction(underlying)
	case *query.Negation:
		return f.generateNegation(underlying)
	case *query.BaseQuery:
		return f.generateBase(underlying)
	case nil:
		return alwaysFalse, nil
	default:
		return alwaysFalse, utils.Should(errors.Errorf("unknown query type: %T", underlying))
	}
}

func (f *Factory) generateDisjunction(q *query.Disjunction) (internalEvaluator, error) {
	return nil, nil
}

func (f *Factory) generateConjunction(q *query.Conjunction) (internalEvaluator, error) {
	return nil, nil
}

func (f *Factory) generateLinkedConjunction(q *query.LinkedConjunction) (internalEvaluator, error) {
	return nil, nil
}

func (f *Factory) generateNegation(q *query.Negation) (internalEvaluator, error) {
	return nil, nil
}

func (f *Factory) generateBase(q *query.BaseQuery) (internalEvaluator, error) {
	fieldPath := f.searchFields.Get(q.Field)
	if len(fieldPath) == 0 {
		return alwaysTrue, nil
	}
	baseEvaluator, err := createBaseEvaluator(fieldPath[len(fieldPath)-1].Type, q.Value)
	if err != nil {
		return nil, errors.Wrapf(err, "generating base query: %v", q)
	}

	pathEvaluator, err := wrapEvaluatorInPath(f.rootType, fieldPath, baseEvaluator)
	if err != nil {
		return nil, errors.Wrapf(err, "generating path traverser: %v", q)
	}
	return pathEvaluator, nil
}
