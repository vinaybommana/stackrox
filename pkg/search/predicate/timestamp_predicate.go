package predicate

import (
	"reflect"

	"github.com/gogo/protobuf/types"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/predicate/basematchers"
)

func createTimestampPredicate(fullPath, value string) (internalPredicate, error) {
	if value == "-" {
		return alwaysFalse, nil
	}

	baseMatcher, err := basematchers.ForTimestamp(value)
	if err != nil {
		return nil, err
	}
	return internalPredicateFunc(func(instance reflect.Value) (*search.Result, bool) {
		instanceTS, _ := instance.Interface().(*types.Timestamp)

		if instanceTS != nil && baseMatcher(instanceTS) {
			return &search.Result{
				Matches: formatSingleMatchf(fullPath, "%d", instanceTS.Seconds),
			}, true
		}
		return nil, false
	}), nil
}
