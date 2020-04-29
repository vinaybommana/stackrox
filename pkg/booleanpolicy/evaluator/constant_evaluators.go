package evaluator

import (
	"reflect"

	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/traverseutil"
)

type alwaysTrueType struct{}

func (alwaysTrueType) Evaluate(path *traverseutil.Path, value reflect.Value) (*Result, bool) {
	return &Result{}, true
}

type alwaysFalseType struct{}

func (alwaysFalseType) Evaluate(path *traverseutil.Path, value reflect.Value) (*Result, bool) {
	return nil, false
}

var (
	alwaysTrue  internalEvaluator = alwaysTrueType{}
	alwaysFalse internalEvaluator = alwaysFalseType{}
)
