package mapeval

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/stringutils"
)

const (
	disjunctionMarker    = ";\t"
	conjunctionMarker    = ",\t"
	shouldNotMatchMarker = "!\t"
)

type kvElement struct {
	key       string
	value     string
	satisfied bool
}

type groupElement struct {
	shouldNotMatch []*kvElement
	shouldMatch    []*kvElement
}

func convertConjunctionPairsToGroupElement(conjunctionPairsStr string) (*groupElement, error) {
	ps := strings.Split(conjunctionPairsStr, conjunctionMarker)
	if len(ps) == 0 {
		return nil, nil
	}

	conjunctionGroup := &groupElement{}
	for _, p := range ps {
		if !strings.Contains(p, "=") {
			return nil, errors.Errorf("Invalid key-value expression: %s", p)
		}

		p, shouldNotMatchQuery := stringutils.MaybeTrimPrefix(p, shouldNotMatchMarker)
		key, value := stringutils.Split2(p, "=")
		ele := &kvElement{value: value, key: key}
		if shouldNotMatchQuery {
			conjunctionGroup.shouldNotMatch = append(conjunctionGroup.shouldNotMatch, ele)
		} else {
			conjunctionGroup.shouldMatch = append(conjunctionGroup.shouldMatch, ele)
		}
	}

	return conjunctionGroup, nil
}

func valueMatchesRequest(req, val string) bool {
	return req == "" || req == val
}

func verifyAgainstCG(gE *groupElement, key, value string) {
	for _, r := range gE.shouldNotMatch {
		r.satisfied = r.satisfied || (valueMatchesRequest(r.key, key) && valueMatchesRequest(r.value, value))
	}

	for _, d := range gE.shouldMatch {
		d.satisfied = d.satisfied || (valueMatchesRequest(d.key, key) && valueMatchesRequest(d.value, value))
	}
}

func matchesCG(gE *groupElement) bool {
	for _, r := range gE.shouldNotMatch {
		if r.satisfied {
			return false
		}
	}
	// All shouldNotMatch requirements failed at this point.

	for _, d := range gE.shouldMatch {
		if !d.satisfied {
			return false
		}
	}
	// Now, all shouldMatch requirements failed at this point, so this map matches this particular conjunction
	// group.
	return true
}

// Matcher returns a matcher for a map against a query string.
func Matcher(value string) (func(*reflect.MapIter) bool, error) {
	// The format for the query is taken to be a disjunction of groups.
	// A group is composed of conjunction of shouldNotMatch and shouldMatch (k,*) (*,v) (k,v) pairs.
	// A shouldMatch pair returns true if it is contained in the map.
	// A shouldNotMatch pair returns true if it is not present in the map.
	// Disjunction is marked by semicolons, Conjunction by commas
	// Should not match groups are preceded by a ! marker, and key value pairs appear as k=v
	// Eg: !a=, b=1; c=2;
	// The above expression is composed of two groups:
	// The first group implies that the map matches if key 'a' is absent, and b=1 is present.
	// The second group implies that the map matches if c=2 is present.
	var disjunctionGroups []*groupElement
	for _, conjunctionPairsStr := range strings.Split(value, disjunctionMarker) {
		cg, err := convertConjunctionPairsToGroupElement(conjunctionPairsStr)
		if err != nil {
			return nil, err
		}

		if cg == nil {
			continue
		}

		disjunctionGroups = append(disjunctionGroups, cg)
	}

	return func(iter *reflect.MapIter) bool {
		for iter.Next() {
			k, v := iter.Key(), iter.Value()
			// Only string type key, value are allowed.
			key, ok := k.Interface().(string)
			if !ok {
				return false
			}

			value, ok := v.Interface().(string)
			if !ok {
				return false
			}

			for _, cg := range disjunctionGroups {
				verifyAgainstCG(cg, key, value)
			}
		}

		for _, cg := range disjunctionGroups {
			// Apply disjunction and return true if any group is true.
			if matchesCG(cg) {
				return true
			}
		}

		return false
	}, nil
}
