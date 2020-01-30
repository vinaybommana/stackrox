package indexer

import (
	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/sync"
)

// Wrapper is an object that wraps keys and values into their indexed id:value pair.
//go:generate mockgen-wrapper
type Wrapper interface {
	Wrap(key []byte, msg proto.Message) (string, interface{})
}

// WrapperRegistry is a registry of all indexers we should use to store messages.
type WrapperRegistry interface {
	RegisterWrapper(prefix []byte, wrapper Wrapper)
	Matches(key []byte) bool

	Wrapper
}

// NewWrapperRegistry returns a new registry for the index.
func NewWrapperRegistry() WrapperRegistry {
	return &wrapperRegistryImpl{
		wrappers: make(map[string]Wrapper),
	}
}

type wrapperRegistryImpl struct {
	lock     sync.RWMutex
	wrappers map[string]Wrapper
}

func (ir *wrapperRegistryImpl) RegisterWrapper(prefix []byte, wrapper Wrapper) {
	concurrency.WithLock(&ir.lock, func() {
		ir.wrappers[string(prefix)] = wrapper
	})
}

func (ir *wrapperRegistryImpl) Matches(key []byte) bool {
	var longestMatch Wrapper
	concurrency.WithRLock(&ir.lock, func() {
		longestMatch = ir.findLongestMatchNoLock(key)
	})
	return longestMatch != nil
}

func (ir *wrapperRegistryImpl) Wrap(key []byte, msg proto.Message) (string, interface{}) {
	var longestMatch Wrapper
	concurrency.WithRLock(&ir.lock, func() {
		longestMatch = ir.findLongestMatchNoLock(key)
	})
	if longestMatch != nil {
		return longestMatch.Wrap(key, msg)
	}
	return "", nil
}

func (ir *wrapperRegistryImpl) findLongestMatchNoLock(key []byte) Wrapper {
	// Need to find the longest matching prefix for a registered index.
	var totalPrefix []byte
	var longestMatch Wrapper
	for currPrefix := badgerhelper.GetPrefix(key); len(currPrefix) > 0; currPrefix = badgerhelper.GetPrefix(badgerhelper.StripPrefix(totalPrefix, key)) {
		totalPrefix = append(totalPrefix, currPrefix...)
		if match, contains := ir.wrappers[string(totalPrefix)]; contains {
			longestMatch = match
		}
	}
	return longestMatch
}
