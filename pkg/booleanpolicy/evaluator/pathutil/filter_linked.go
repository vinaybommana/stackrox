package pathutil

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/utils"
)

type tree struct {
	children map[stepMapKey]*tree
}

func newTree() *tree {
	return &tree{
		children: make(map[stepMapKey]*tree),
	}
}

func (t *tree) addPath(steps []step) {
	if len(steps) == 0 {
		return
	}
	firstStep, remainingSteps := steps[0], steps[1:]
	key := firstStep.mapKey()
	subTree := t.children[key]
	if subTree == nil {
		subTree = newTree()
		t.children[key] = subTree
	}
	subTree.addPath(remainingSteps)
}

// treeFromPaths generates a tree from the given paths.
// Callers must ensure that:
// a) there is at least one path
// b) all the paths are of the same length
func treeFromPaths(pathHolders []PathHolder) *tree {
	t := newTree()
	for _, pathHolder := range pathHolders {
		path := pathHolder.GetPath()
		if len(path.steps) == 0 {
			utils.Should(errors.Errorf("empty path from search (paths: %v)", pathHolders))
			continue
		}
		t.addPath(path.steps[:len(path.steps)-1])
	}
	return t
}

func (t *tree) merge(other *tree) {
	for key, child := range t.children {
		otherChild, inOther := other.children[key]
		if inOther {
			child.merge(otherChild)
			continue
		}
		// For integer values, which represent an array index, we must drop unless the value is in both.
		if _, isInt := key.(int); isInt {
			delete(t.children, key)
		}
	}
	for key, child := range other.children {
		if _, inT := t.children[key]; inT {
			// This key has been considered already in the above loop.
			continue
		}
		// Don't merge integer keys unless they're in both.
		if _, isInt := key.(int); isInt {
			continue
		}
		// Copy over the child.
		t.children[key] = child
	}
}

func (t *tree) filterToMatchingPaths(field string, pathHolders []PathHolder) []PathHolder {
	filtered := pathHolders[:0]
	for _, pathHolder := range pathHolders {
		path := pathHolder.GetPath()
		// This is an invalid path, should never happen. The panic will be caught and softened to a utils.Should
		// by the caller.
		if len(path.steps) == 0 {
			panic(fmt.Sprintf("invalid: got empty path for field %s", field))
		}
		if t.containsSteps(path.steps[:len(path.steps)-1]) {
			filtered = append(filtered, pathHolder)
		}
	}
	return filtered
}

func (t *tree) containsSteps(steps []step) bool {
	// Base case
	if len(steps) == 0 {
		return true
	}
	firstStep := steps[0]
	child := t.children[firstStep.mapKey()]
	if child == nil {
		return false
	}
	return child.containsSteps(steps[1:])
}

// A PathHolder is any object containing a path.
type PathHolder interface {
	GetPath() *Path
}

// FilterPathsToLinkedMatches filters the given fieldsToPaths to just the linked matches.
// The best way to understand the purpose of this function is to look at the unit tests.
func FilterPathsToLinkedMatches(fieldsToPaths map[string][]PathHolder) (fieldsToFilteredPaths map[string][]PathHolder, matched bool, err error) {
	// For convenience, the internal functions here signal errors by panic-ing, but we catch the panic here
	// so that clients outside the package just receive an error.
	// Panics will only happen with invalid inputs, which is always a programming error.
	defer func() {
		if r := recover(); r != nil {
			err = utils.Should(errors.Errorf("invalid input: %v", r))
		}
	}()
	t := newTree()
	for _, paths := range fieldsToPaths {
		t.merge(treeFromPaths(paths))
	}

	fieldsToFilteredPaths = make(map[string][]PathHolder)
	for field, paths := range fieldsToPaths {
		filtered := t.filterToMatchingPaths(field, paths)
		if len(filtered) == 0 {
			return nil, false, nil
		}
		fieldsToFilteredPaths[field] = filtered
	}

	return fieldsToFilteredPaths, true, nil
}
