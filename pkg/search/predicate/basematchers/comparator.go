package basematchers

import (
	"fmt"
	"strings"
)

// Produce a predicate for the given numerical or date time query.
const (
	lessThanOrEqualTo    = "<="
	greaterThanOrEqualTo = ">="
	lessThan             = "<"
	greaterThan          = ">"
)

func intComparator(cmp string) (func(a, b int64) bool, error) {
	switch cmp {
	case lessThanOrEqualTo:
		return func(a, b int64) bool { return a <= b }, nil
	case greaterThanOrEqualTo:
		return func(a, b int64) bool { return a >= b }, nil
	case lessThan:
		return func(a, b int64) bool { return a < b }, nil
	case greaterThan:
		return func(a, b int64) bool { return a > b }, nil
	case "":
		return func(a, b int64) bool { return a == b }, nil
	default:
		return nil, fmt.Errorf("unrecognized comparator: %s", cmp)
	}
}

func uintComparator(cmp string) (func(a, b uint64) bool, error) {
	switch cmp {
	case lessThanOrEqualTo:
		return func(a, b uint64) bool { return a <= b }, nil
	case greaterThanOrEqualTo:
		return func(a, b uint64) bool { return a >= b }, nil
	case lessThan:
		return func(a, b uint64) bool { return a < b }, nil
	case greaterThan:
		return func(a, b uint64) bool { return a > b }, nil
	case "":
		return func(a, b uint64) bool { return a == b }, nil
	default:
		return nil, fmt.Errorf("unrecognized comparator: %s", cmp)
	}
}

func floatComparator(cmp string) (func(a, b float64) bool, error) {
	switch cmp {
	case lessThanOrEqualTo:
		return func(a, b float64) bool { return a <= b }, nil
	case greaterThanOrEqualTo:
		return func(a, b float64) bool { return a >= b }, nil
	case lessThan:
		return func(a, b float64) bool { return a < b }, nil
	case greaterThan:
		return func(a, b float64) bool { return a > b }, nil
	case "":
		return func(a, b float64) bool { return a == b }, nil
	default:
		return nil, fmt.Errorf("unrecognized comparator: %s", cmp)
	}
}

func parseNumericPrefix(value string) (prefix string, trimmedValue string) {
	// The order which these checks are executed must be maintained.
	// If we for instance look for "<" before "<=", we will never find "<=" because "<" will be found as its prefix.
	for _, prefix := range []string{lessThanOrEqualTo, greaterThanOrEqualTo, lessThan, greaterThan} {
		if strings.HasPrefix(value, prefix) {
			return prefix, value[len(prefix):]
		}
	}
	return "", value
}
