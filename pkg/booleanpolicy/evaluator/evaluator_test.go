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
}

type SecondNested struct {
	SecondNestedValA string `search:"SecondA" protobuf:"blah"`
	SecondNestedValB string `search:"SecondB" protobuf:"blah"`
}

func TestSimpleBase(t *testing.T) {
	factory := NewFactory((*TopLevel)(nil))

	qTopLevelAHappy := query.NewBase("TopLevelA", "happy")
	qNestedAHappy := query.NewBase("A", "happy")
	qSecondNestedAHappy := query.NewBase("SecondA", "r/.*ppy")

	for _, testCase := range []struct {
		desc           string
		q              query.Query
		obj            *TopLevel
		expectedResult *Result
	}{
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
			expectedResult: resultWithSingleMatch(traverseutil.PathFromSteps(t, "ValA"), "happy"),
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
			expectedResult: resultWithSingleMatch(traverseutil.PathFromSteps(t, "NestedSlice", 0, "NestedValA"), "happy"),
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
			expectedResult: &Result{Matches: []Match{
				{Path: traverseutil.PathFromSteps(t, "NestedSlice", 0, "SecondNestedSlice", 1, "SecondNestedValA"), Value: "blaappy"},
				{Path: traverseutil.PathFromSteps(t, "NestedSlice", 1, "SecondNestedSlice", 0, "SecondNestedValA"), Value: "happy"},
			},
			},
		},
	} {
		c := testCase
		t.Run(c.desc, func(t *testing.T) {
			evaluator, err := factory.GenerateEvaluator(c.q)
			require.NoError(t, err)
			res, matched := evaluator.Evaluate(c.obj)
			assert.Equal(t, c.expectedResult != nil, matched)
			if c.expectedResult != nil {
				require.NotNil(t, res)
				assert.ElementsMatch(t, c.expectedResult.Matches, res.Matches)
			}
		})
	}
}
