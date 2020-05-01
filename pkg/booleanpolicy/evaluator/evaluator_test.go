package evaluator

import (
	"testing"

	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/traverseutil"
	"github.com/stackrox/rox/pkg/booleanpolicy/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TopLevel struct {
	ValA        string   `search:"TopLevelA" protobuf:"blah"`
	NestedSlice []Nested `protobuf:"blah"`
}

type Nested struct {
	NestedValA        string          `search:"A" protobuf:"blah"`
	NestedValB        string          `search:"B" protobuf:"blah"`
	SecondNestedSlice []*SecondNested `protobuf:"blah"`
	ThirdNestedSlice  []ThirdNested   `protobuf:"blah"`
}

type SecondNested struct {
	SecondNestedValA string `search:"SecondA" protobuf:"blah"`
	SecondNestedValB string `search:"SecondB" protobuf:"blah"`
}

type ThirdNested struct {
	ThirdNestedValA string   `search:"ThirdA" protobuf:"blah"`
	ThirdNestedValB []string `search:"ThirdB" protobuf:"blah"`
}

var (
	factoryInstance = NewFactory((*TopLevel)(nil))
)

type testCase struct {
	desc           string
	q              *query.Query
	obj            *TopLevel
	expectedResult *Result
}

func runTestCases(t *testing.T, testCases []testCase) {
	for _, testCase := range testCases {
		c := testCase
		t.Run(c.desc, func(t *testing.T) {
			evaluator, err := factoryInstance.GenerateEvaluator(c.q)
			require.NoError(t, err)
			res, matched := evaluator.Evaluate(c.obj)
			assert.Equal(t, c.expectedResult != nil, matched)
			if c.expectedResult != nil {
				require.NotNil(t, res)
				assert.Equal(t, c.expectedResult.Matches, res.Matches)
			}
		})

	}
}

func TestSimpleBase(t *testing.T) {
	qTopLevelAHappy := query.SimpleMatchFieldQuery("TopLevelA", "happy")
	qNestedAHappy := query.SimpleMatchFieldQuery("A", "happy")
	qSecondNestedAHappy := query.SimpleMatchFieldQuery("SecondA", "r/.*ppy")

	runTestCases(t, []testCase{
		{
			desc: "simple one for top level, doesn't pass",
			q:    qTopLevelAHappy,
			obj: &TopLevel{
				ValA: "whatever",
				NestedSlice: []Nested{
					{NestedValA: "blah"},
					{NestedValA: "something else", SecondNestedSlice: []*SecondNested{
						{SecondNestedValA: "happy"},
					}},
				},
			},
		},
		{
			desc: "simple one for top level, passes",
			q:    qTopLevelAHappy,
			obj: &TopLevel{
				ValA: "happy",
				NestedSlice: []Nested{
					{NestedValA: "blah"},
					{NestedValA: "something else", SecondNestedSlice: []*SecondNested{
						{SecondNestedValA: "happy"},
					}},
				},
			},
			expectedResult: resultWithSingleMatch("TopLevelA", traverseutil.PathFromSteps(t, "ValA"), "happy"),
		},
		{
			desc: "simple one for first level nested, doesn't pass",
			q:    qNestedAHappy,
			obj: &TopLevel{
				ValA: "happy",
				NestedSlice: []Nested{
					{NestedValA: "blah"},
					{NestedValA: "something else", SecondNestedSlice: []*SecondNested{
						{SecondNestedValA: "happy"},
					}},
				},
			},
		},
		{
			desc: "simple one for first level nested, passes",
			q:    qNestedAHappy,
			obj: &TopLevel{
				ValA: "happy",
				NestedSlice: []Nested{
					{NestedValA: "happy"},
					{NestedValA: "something else", SecondNestedSlice: []*SecondNested{
						{SecondNestedValA: "happiest"},
					}},
				},
			},
			expectedResult: resultWithSingleMatch("A", traverseutil.PathFromSteps(t, "NestedSlice", 0, "NestedValA"), "happy"),
		},
		{
			desc: "simple one for second level nested, doesn't pass",
			q:    qSecondNestedAHappy,
			obj: &TopLevel{
				ValA: "happy",
				NestedSlice: []Nested{
					{NestedValA: "happy"},
					{NestedValA: "something else", SecondNestedSlice: []*SecondNested{
						{SecondNestedValA: "happiest"},
					}},
				},
			},
		},
		{
			desc: "simple one for second level nested, passes",
			q:    qSecondNestedAHappy,
			obj: &TopLevel{
				ValA: "happy",
				NestedSlice: []Nested{
					{NestedValA: "happy", SecondNestedSlice: []*SecondNested{
						{SecondNestedValA: "blah"},
						{SecondNestedValA: "blaappy"},
					}},
					{NestedValA: "something else", SecondNestedSlice: []*SecondNested{
						{SecondNestedValA: "happy"},
					}},
				},
			},
			expectedResult: &Result{Matches: map[string][]Match{
				"SecondA": {
					{Path: traverseutil.PathFromSteps(t, "NestedSlice", 0, "SecondNestedSlice", 1, "SecondNestedValA"), Values: []string{"blaappy"}},
					{Path: traverseutil.PathFromSteps(t, "NestedSlice", 1, "SecondNestedSlice", 0, "SecondNestedValA"), Values: []string{"happy"}},
				},
			},
			},
		},
	})
}

func TestLinked(t *testing.T) {
	runTestCases(t, []testCase{
		{
			desc: "linked, first level of nesting, should match",
			obj: &TopLevel{
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "A", Values: []string{"A1"}},
					{Field: "B", Values: []string{"B1"}},
				},
			},
			expectedResult: &Result{Matches: map[string][]Match{
				"A": {{Path: traverseutil.PathFromSteps(t, "NestedSlice", 1, "NestedValA"), Values: []string{"A1"}}},
				"B": {{Path: traverseutil.PathFromSteps(t, "NestedSlice", 1, "NestedValB"), Values: []string{"B1"}}},
			}},
		},
		{
			desc: "linked, first level of nesting, should not match",
			obj: &TopLevel{
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "A", Values: []string{"A0"}},
					{Field: "B", Values: []string{"B1"}},
				},
			},
		},
		{
			desc: "linked, multilevel, should match",
			obj: &TopLevel{
				ValA: "TopLevelValA",
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "TopLevelA", Values: []string{"TopLevelValA"}},
					{Field: "A", Values: []string{"A1"}},
					{Field: "B", Values: []string{"B1"}},
				},
			},
			expectedResult: &Result{Matches: map[string][]Match{
				"TopLevelA": {{Path: traverseutil.PathFromSteps(t, "ValA"), Values: []string{"TopLevelValA"}}},
				"A":         {{Path: traverseutil.PathFromSteps(t, "NestedSlice", 1, "NestedValA"), Values: []string{"A1"}}},
				"B":         {{Path: traverseutil.PathFromSteps(t, "NestedSlice", 1, "NestedValB"), Values: []string{"B1"}}},
			}},
		},
		{
			desc: "linked, multilevel, top doesn't match",
			obj: &TopLevel{
				ValA: "TopLevelValA",
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "TopLevelA", Values: []string{"NONEXISTENT"}},
					{Field: "A", Values: []string{"A1"}},
					{Field: "B", Values: []string{"B1"}},
				},
			},
		},
		{
			desc: "linked, multilevel, bottom doesn't match",
			obj: &TopLevel{
				ValA: "TopLevelValA",
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "TopLevelA", Values: []string{"TopLevelValA"}},
					{Field: "A", Values: []string{"A0"}},
					{Field: "B", Values: []string{"B1"}},
				},
			},
		},
	})
}

func TestCompound(t *testing.T) {
	runTestCases(t, []testCase{
		{
			desc: "simple compound query, OR, matches",
			obj: &TopLevel{
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "A", Values: []string{"A0", "A1"}, Operator: query.Or},
				},
			},
			expectedResult: &Result{Matches: map[string][]Match{
				"A": {
					{Path: traverseutil.PathFromSteps(t, "NestedSlice", 0, "NestedValA"), Values: []string{"A0"}},
					{Path: traverseutil.PathFromSteps(t, "NestedSlice", 1, "NestedValA"), Values: []string{"A1"}},
				},
			}},
		},
		{
			desc: "simple compound query, OR, does not match",
			obj: &TopLevel{
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "A", Values: []string{"A2", "A3"}, Operator: query.Or},
				},
			},
		},
		{
			desc: "simple compound query, AND, does not match",
			obj: &TopLevel{
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "A", Values: []string{"A0", "A1"}, Operator: query.And},
				},
			},
		},
		{
			desc: "simple compound query, AND, matches",
			obj: &TopLevel{
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "A", Values: []string{"r/A.*", "r/.*1"}, Operator: query.And},
				},
			},
			expectedResult: &Result{Matches: map[string][]Match{
				"A": {{Path: traverseutil.PathFromSteps(t, "NestedSlice", 1, "NestedValA"), Values: []string{"A1"}}},
			}},
		},
		{
			desc: "compound query, OR, negated, matches",
			obj: &TopLevel{
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "A", Values: []string{"A2", "A1"}, Operator: query.Or, Negate: true},
				},
			},
			expectedResult: &Result{Matches: map[string][]Match{
				"A": {
					{Path: traverseutil.PathFromSteps(t, "NestedSlice", 0, "NestedValA"), Values: []string{"A0"}},
				},
			}},
		},
		{
			desc: "compound query, OR, negated, does not match",
			obj: &TopLevel{
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "A", Values: []string{"A0", "A1"}, Operator: query.Or, Negate: true},
				},
			},
		},
		{
			desc: "compound query, AND, negated, does not match",
			obj: &TopLevel{
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "A", Values: []string{`r/A.*`, `r/.*\d`}, Operator: query.And, Negate: true},
				},
			},
		},
		{
			desc: "simple compound query, AND, negated, matches",
			obj: &TopLevel{
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "A", Values: []string{"r/A.*", "r/.*1"}, Operator: query.And, Negate: true},
				},
			},
			expectedResult: &Result{Matches: map[string][]Match{
				"A": {{Path: traverseutil.PathFromSteps(t, "NestedSlice", 0, "NestedValA"), Values: []string{"A0"}}},
			}},
		},
	})

}
