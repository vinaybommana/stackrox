package pathutil

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type nestedIgnored struct {
	D string `search:"D"`
}

type nested struct {
	B string `search:"B"`
}

type topLevel struct {
	A             int `search:"A"`
	Nested        []nested
	NestedIgnored nestedIgnored `search:"-"`
}

type objWithInt struct {
	AugmentedVal int `search:"AugmentedInt"`
}

type objWithString struct {
	AugmentedVal string `search:"AugmentedStr"`
}

func TestAugmentedObjMeta(t *testing.T) {
	for _, testCase := range []struct {
		desc             string
		objMeta          *AugmentedObjMeta
		expectedFieldMap map[string]MetaPath
		shouldErr        bool
	}{
		{
			desc:    "plain object, unaugmented",
			objMeta: NewAugmentedObjMeta((*topLevel)(nil)),
			expectedFieldMap: map[string]MetaPath{
				"a": {{FieldName: "A", Type: reflect.TypeOf(0), StructFieldIndex: []int{0}}},
				"b": {
					{FieldName: "Nested", Type: reflect.TypeOf(([]nested)(nil)), StructFieldIndex: []int{1}},
					{FieldName: "B", Type: reflect.TypeOf(""), StructFieldIndex: []int{0}},
				},
			},
		},
		{
			desc: "augmented at top level",
			objMeta: NewAugmentedObjMeta((*topLevel)(nil)).
				AddPlainObjectAt([]string{"IntObj"}, (*objWithInt)(nil)),
			expectedFieldMap: map[string]MetaPath{
				"a": {{FieldName: "A", Type: reflect.TypeOf(0), StructFieldIndex: []int{0}}},
				"b": {
					{FieldName: "Nested", Type: reflect.TypeOf(([]nested)(nil)), StructFieldIndex: []int{1}},
					{FieldName: "B", Type: reflect.TypeOf(""), StructFieldIndex: []int{0}},
				},
				"augmentedint": {
					{FieldName: "IntObj", Type: reflect.TypeOf((*objWithInt)(nil))},
					{FieldName: "AugmentedVal", Type: reflect.TypeOf(0), StructFieldIndex: []int{0}},
				},
			},
		},
		{
			desc: "augmented inside nested",
			objMeta: NewAugmentedObjMeta((*topLevel)(nil)).
				AddPlainObjectAt([]string{"Nested", "IntObj"}, (*objWithInt)(nil)),
			expectedFieldMap: map[string]MetaPath{
				"a": {{FieldName: "A", Type: reflect.TypeOf(0), StructFieldIndex: []int{0}}},
				"b": {
					{FieldName: "Nested", Type: reflect.TypeOf(([]nested)(nil)), StructFieldIndex: []int{1}},
					{FieldName: "B", Type: reflect.TypeOf(""), StructFieldIndex: []int{0}},
				},
				"augmentedint": {
					{FieldName: "Nested", Type: reflect.TypeOf(([]nested)(nil)), StructFieldIndex: []int{1}},
					{FieldName: "IntObj", Type: reflect.TypeOf((*objWithInt)(nil))},
					{FieldName: "AugmentedVal", Type: reflect.TypeOf(0), StructFieldIndex: []int{0}},
				},
			},
		},
		{
			desc: "multiple augments",
			objMeta: NewAugmentedObjMeta((*topLevel)(nil)).
				AddPlainObjectAt([]string{"Nested", "IntObj"}, (*objWithInt)(nil)).
				AddPlainObjectAt([]string{"StringObj"}, (*objWithString)(nil)),
			expectedFieldMap: map[string]MetaPath{
				"a": {{FieldName: "A", Type: reflect.TypeOf(0), StructFieldIndex: []int{0}}},
				"b": {
					{FieldName: "Nested", Type: reflect.TypeOf(([]nested)(nil)), StructFieldIndex: []int{1}},
					{FieldName: "B", Type: reflect.TypeOf(""), StructFieldIndex: []int{0}},
				},
				"augmentedint": {
					{FieldName: "Nested", Type: reflect.TypeOf(([]nested)(nil)), StructFieldIndex: []int{1}},
					{FieldName: "IntObj", Type: reflect.TypeOf((*objWithInt)(nil))},
					{FieldName: "AugmentedVal", Type: reflect.TypeOf(0), StructFieldIndex: []int{0}},
				},
				"augmentedstr": {
					{FieldName: "StringObj", Type: reflect.TypeOf((*objWithString)(nil))},
					{FieldName: "AugmentedVal", Type: reflect.TypeOf(""), StructFieldIndex: []int{0}},
				},
			},
		},
		{
			desc: "augments within augments",
			objMeta: NewAugmentedObjMeta((*topLevel)(nil)).
				AddAugmentedObjectAt(
					[]string{"Nested", "IntObj"},
					NewAugmentedObjMeta((*objWithInt)(nil)).
						AddPlainObjectAt([]string{"StringObj"}, (*objWithString)(nil)),
				),
			expectedFieldMap: map[string]MetaPath{
				"a": {{FieldName: "A", Type: reflect.TypeOf(0), StructFieldIndex: []int{0}}},
				"b": {
					{FieldName: "Nested", Type: reflect.TypeOf(([]nested)(nil)), StructFieldIndex: []int{1}},
					{FieldName: "B", Type: reflect.TypeOf(""), StructFieldIndex: []int{0}},
				},
				"augmentedint": {
					{FieldName: "Nested", Type: reflect.TypeOf(([]nested)(nil)), StructFieldIndex: []int{1}},
					{FieldName: "IntObj", Type: reflect.TypeOf((*objWithInt)(nil))},
					{FieldName: "AugmentedVal", Type: reflect.TypeOf(0), StructFieldIndex: []int{0}},
				},
				"augmentedstr": {
					{FieldName: "Nested", Type: reflect.TypeOf(([]nested)(nil)), StructFieldIndex: []int{1}},
					{FieldName: "IntObj", Type: reflect.TypeOf((*objWithInt)(nil))},
					{FieldName: "StringObj", Type: reflect.TypeOf((*objWithString)(nil))},
					{FieldName: "AugmentedVal", Type: reflect.TypeOf(""), StructFieldIndex: []int{0}},
				},
			},
		},
		{
			desc: "an augment with a name clash, should replace the original object",
			objMeta: NewAugmentedObjMeta((*topLevel)(nil)).
				AddPlainObjectAt([]string{"Nested"}, (*objWithInt)(nil)),
			expectedFieldMap: map[string]MetaPath{
				"a": {{FieldName: "A", Type: reflect.TypeOf(0), StructFieldIndex: []int{0}}},
				"augmentedint": {
					{FieldName: "Nested", Type: reflect.TypeOf((*objWithInt)(nil))},
					{FieldName: "AugmentedVal", Type: reflect.TypeOf(0), StructFieldIndex: []int{0}},
				},
			},
		},
		{
			desc: "invalid augment, invalid path",
			objMeta: NewAugmentedObjMeta((*topLevel)(nil)).
				AddPlainObjectAt([]string{"NonExistent", "IntObj"}, (*objWithInt)(nil)),
			shouldErr: true,
		},
	} {
		c := testCase
		t.Run(c.desc, func(t *testing.T) {
			out, err := c.objMeta.MapSearchTagsToPaths()
			if c.shouldErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, c.expectedFieldMap, out.underlying)
		})
	}
}
