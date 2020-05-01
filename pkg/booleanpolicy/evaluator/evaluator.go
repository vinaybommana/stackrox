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
func (f *Factory) GenerateEvaluator(q *query.Query) (Evaluator, error) {
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

func (f *Factory) generateInternalEvaluator(q *query.Query) (internalEvaluator, error) {
	// The field queries are implicitly a linked conjunction. This means that all field queries must match,
	// AND that their matches must be in the same object.
	// The notion of linking is a bit complicated -- the easiest way to get a sense of what it entails is to look
	// at the test cases in TestLinked.
	fieldEvaluators := make([]internalEvaluator, 0, len(q.FieldQueries))
	for _, fq := range q.FieldQueries {
		eval, err := f.generateInternalEvaluatorForFieldQuery(fq)
		if err != nil {
			return nil, errors.Wrapf(err, "compiling field query: %v", fq)
		}

		// If one of them is alwaysFalse, then our conjunction is alwaysFalse.
		if eval == alwaysFalse {
			return alwaysFalse, nil
		}
		if eval != alwaysTrue {
			fieldEvaluators = append(fieldEvaluators, eval)
		}
	}
	switch len(fieldEvaluators) {
	case 0:
		return alwaysTrue, nil
	case 1:
		// Simplify the case where there's just one.
		return fieldEvaluators[0], nil
	default:
		return internalEvaluatorFunc(func(path *traverseutil.Path, value reflect.Value) (*Result, bool) {
			fieldsToPaths := make(map[string][]traverseutil.PathHolder)
			for _, fieldEval := range fieldEvaluators {
				result, matches := fieldEval.Evaluate(path, value)
				if !matches {
					return nil, false
				}
				for field, matches := range result.Matches {
					for _, match := range matches {
						fieldsToPaths[field] = append(fieldsToPaths[field], match)
					}
				}
			}
			filteredFieldsToPaths, matched, err := traverseutil.FilterPathsToLinkedMatches(fieldsToPaths)
			if err != nil {
				utils.Should(errors.Wrap(err, "filtering paths to linked matches"))
				return nil, false
			}
			if !matched {
				return nil, false
			}
			r := newResult()
			for field, matches := range filteredFieldsToPaths {
				for _, match := range matches {
					r.Matches[field] = append(r.Matches[field], match.(Match))
				}
			}
			return r, true
		}), nil
	}
}

// generateInternalEvaluatorForFieldQuery generates an internal evaluator for a specific field query.
func (f *Factory) generateInternalEvaluatorForFieldQuery(q *query.FieldQuery) (internalEvaluator, error) {
	fieldPath := f.searchFields.Get(q.Field)
	// This happens for the case where we're in a factory for images, but are querying a field that exists only
	// on deployments. It's irrelevant to the image query if that's the case.
	if len(fieldPath) == 0 {
		return alwaysTrue, nil
	}

	baseType := fieldPath[len(fieldPath)-1].Type
	baseEvaluator, err := createBaseEvaluator(q.Field, baseType, q.Values, q.Negate, q.Operator)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid query %v", q)
	}

	pathEvaluator, err := wrapEvaluatorInPath(f.rootType, fieldPath, baseEvaluator)
	if err != nil {
		return nil, errors.Wrapf(err, "generating path traverser: %v", q)
	}
	return pathEvaluator, nil
}
