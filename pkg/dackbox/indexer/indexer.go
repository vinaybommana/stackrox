package indexer

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/sync"
)

// Indexer is an object that is able to index proto messages.
//go:generate mockgen-wrapper
type Indexer interface {
	Index(key []byte, msg proto.Message) error
	Delete(key []byte) error
}

// IndexRegistry is a registry of all indexers we should use to store messages.
type IndexRegistry interface {
	RegisterIndex(prefix []byte, index Indexer)
	Matches(key []byte) bool

	Indexer
}

// NewIndexRegistry returns a new registry for the index.
func NewIndexRegistry() IndexRegistry {
	return &indexRegistryImpl{
		indexes: make(map[string]Indexer),
	}
}

type indexRegistryImpl struct {
	lock    sync.RWMutex
	indexes map[string]Indexer
}

func (ir *indexRegistryImpl) RegisterIndex(prefix []byte, index Indexer) {
	ir.lock.Lock()
	defer ir.lock.Unlock()

	ir.indexes[string(prefix)] = index
}

func (ir *indexRegistryImpl) Matches(key []byte) bool {
	ir.lock.RLock()
	defer ir.lock.RUnlock()

	_, longestMatch := ir.findLongestMatchNoLock(key)
	return longestMatch != nil
}

func (ir *indexRegistryImpl) Index(key []byte, msg proto.Message) error {
	ir.lock.RLock()
	defer ir.lock.RUnlock()

	longestPrefix, longestMatch := ir.findLongestMatchNoLock(key)
	if longestMatch != nil {
		return longestMatch.Index(badgerhelper.StripBucket(longestPrefix, key), msg)
	}
	return errors.Errorf("cannot index: no indexer registered matched input key %s", string(key))
}

func (ir *indexRegistryImpl) Delete(key []byte) error {
	ir.lock.RLock()
	defer ir.lock.RUnlock()

	longestPrefix, longestMatch := ir.findLongestMatchNoLock(key)
	if longestMatch != nil {
		return longestMatch.Delete(badgerhelper.StripBucket(longestPrefix, key))
	}
	return errors.Errorf("cannot delete: no indexer registered matched input key %s", string(key))
}

func (ir *indexRegistryImpl) findLongestMatchNoLock(key []byte) ([]byte, Indexer) {
	// Need to find the longest matching prefix for a registered index.
	var totalPrefix []byte
	var longestPrefix []byte
	var longestMatch Indexer
	for currPrefix := badgerhelper.GetPrefix(key); len(currPrefix) > 0; currPrefix = badgerhelper.GetPrefix(badgerhelper.StripPrefix(totalPrefix, key)) {
		totalPrefix = append(totalPrefix, currPrefix...)
		if match, contains := ir.indexes[string(totalPrefix)]; contains {
			longestPrefix = totalPrefix
			longestMatch = match
		}
	}
	return longestPrefix, longestMatch
}
