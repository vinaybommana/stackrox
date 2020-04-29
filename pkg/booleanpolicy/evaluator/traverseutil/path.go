package traverseutil

import (
	"github.com/stackrox/rox/pkg/pointers"
)

type step struct {
	field string
	index *int
}

// A Path represents a list of steps taken to traverse an object.
// This includes struct field indirections and array indexing.
// Paths are copy-on-write.
type Path struct {
	steps []step
}

func (p *Path) cloneAndAddStep(newStep step) *Path {
	newPath := &Path{steps: make([]step, len(p.steps)+1)}
	copy(newPath.steps, p.steps)
	newPath.steps[len(p.steps)] = newStep
	return newPath
}

// WithFieldTraversed returns a copy of path that adds a new step
// that traverses a struct field.
func (p *Path) WithFieldTraversed(fieldName string) *Path {
	return p.cloneAndAddStep(step{field: fieldName})
}

// WithSliceIndexed returns a copy of path that adds a new step
// that indexes into a slice.
func (p *Path) WithSliceIndexed(index int) *Path {
	return p.cloneAndAddStep(step{index: pointers.Int(index)})
}
