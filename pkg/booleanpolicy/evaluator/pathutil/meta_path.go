package pathutil

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// A MetaStep represents a step in a MetaPath. It is either a struct field
// traversal (through StructFieldIndex), or a leap through an augmented object.
type MetaStep struct {
	Type             reflect.Type
	FieldName        string
	StructFieldIndex []int // This is reflect.StructField.Index, which is an efficient way to index into a struct.
}

// A MetaPath represents a path on types
// (ie, a MetaPath can be thought of the abstract version of a Path).
// Whereas a Path operates on an _instance_ of an object,
// a MetaPath operates on the type.
type MetaPath []MetaStep

// FieldToMetaPathMap helps store and retrieve meta paths given a field tag.
type FieldToMetaPathMap struct {
	underlying map[string]MetaPath
}

func (m *FieldToMetaPathMap) maybeAdd(tag string, metaPath MetaPath) {
	_ = m.add(tag, metaPath)
}

func (m *FieldToMetaPathMap) add(tag string, metaPath MetaPath) error {
	lowerTag := strings.ToLower(tag)
	if existingPath, exists := m.underlying[lowerTag]; exists {
		return errors.Errorf("duplicate search tag detected: %s (clashing paths: %v/%v)", tag, existingPath, metaPath)
	}
	m.underlying[lowerTag] = metaPath
	return nil
}

// Get returns the MetaPath for the given tag, and a bool indicates whether it exists.
func (m *FieldToMetaPathMap) Get(tag string) (MetaPath, bool) {
	metaPath, found := m.underlying[strings.ToLower(tag)]
	return metaPath, found
}
