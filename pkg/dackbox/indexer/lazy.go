package indexer

import (
	"strings"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/dackbox/utils/queue"
	"github.com/stackrox/rox/pkg/logging"
)

const (
	procInterval = 500 * time.Millisecond
	maxBatchSize = 500
)

var (
	log = logging.LoggerForModule()
)

// Acker is a function we call on keys that have been processed.
type Acker func(keys ...[]byte) error

// Lazy represents an interface for lazily indexing values that have been written to DackBox.
type Lazy interface {
	Mark([]byte, proto.Message)

	Start()
	Stop()
}

// NewLazy returns a new instance of a lazy indexer that reads in the values to index from the toIndex queue, indexes
// them with the given indexer, then acks indexed values with the given acker.
func NewLazy(toIndex queue.WaitableQueue, wrapper Wrapper, index bleve.Index, acker Acker) Lazy {
	return &lazyImpl{
		wrapper:    wrapper,
		index:      index,
		acker:      acker,
		toIndex:    toIndex,
		stopSignal: concurrency.NewSignal(),
	}
}

type lazyImpl struct {
	wrapper Wrapper
	index   bleve.Index
	acker   Acker
	toIndex queue.WaitableQueue

	stopSignal concurrency.Signal
}

func (li *lazyImpl) Mark(key []byte, value proto.Message) {
	li.toIndex.Push(key, value)
}

func (li *lazyImpl) Start() {
	go li.runIndexing()
}

func (li *lazyImpl) Stop() {
	li.stopSignal.Signal()
}

// No need for control logic since we always want this running with an instance of DackBox that uses lazy indexing.
func (li *lazyImpl) runIndexing() {
	ticker := time.NewTicker(procInterval)
	defer ticker.Stop()

	valuesToIndex := make(map[string]interface{})
	keysToAck := make([][]byte, 0, maxBatchSize)
	for {
		select {
		case <-li.stopSignal.Done():
			return

		// Don't wait more than the interval to index.
		case <-ticker.C:
			li.flush(keysToAck, valuesToIndex)
			keysToAck = keysToAck[:0]
			valuesToIndex = make(map[string]interface{})

		// Collect items from the queue to index.
		case <-li.toIndex.NotEmpty().Done():
			key, value := li.toIndex.Pop()
			if key == nil {
				continue
			}

			indexedKey, indexedValue := li.wrapper.Wrap(key, value)
			if indexedKey == "" {
				log.Errorf("no wrapper registered for key: %s", string(key))
				continue
			}
			keysToAck = append(keysToAck, key)
			valuesToIndex[indexedKey] = indexedValue

			// Don't ack more than the max at a time.
			if len(keysToAck) >= maxBatchSize || len(valuesToIndex) >= maxBatchSize {
				li.flush(keysToAck, valuesToIndex)
				keysToAck = keysToAck[:0]
				valuesToIndex = make(map[string]interface{})
			}
		}
	}
}

func (li *lazyImpl) flush(keysToAck [][]byte, valuesToIndex map[string]interface{}) {
	li.indexItems(valuesToIndex)
	li.ackKeys(keysToAck)
}

func (li *lazyImpl) indexItems(itemsToIndex map[string]interface{}) {
	batch := li.index.NewBatch()
	for key, value := range itemsToIndex {
		if value != nil {
			if err := batch.Index(key, value); err != nil {
				log.Errorf("unable to index item: %s, %v", key, err)
			}
		} else {
			batch.Delete(key)
		}
	}
	err := li.index.Batch(batch)
	if err != nil {
		log.Errorf("unable to index batch of items: %v", err)
	}
}

func (li *lazyImpl) ackKeys(keysToAck [][]byte) {
	err := li.acker(keysToAck...)
	if err != nil {
		log.Errorf("unable to ack keys: %s, %v", printableKeys(keysToAck), err)
	}
}

// Helper for printing key values.
type printableKeys [][]byte

func (pk printableKeys) String() string {
	keys := make([]string, 0, len(pk))
	for _, key := range keys {
		keys = append(keys, string(key))
	}
	return strings.Join(keys, ", ")
}
