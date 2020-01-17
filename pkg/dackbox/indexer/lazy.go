package indexer

import (
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/dackbox/utils/queue"
	"github.com/stackrox/rox/pkg/logging"
)

const (
	ackInterval = 10 * time.Second
	maxToAck    = 1000
)

var (
	log = logging.LoggerForModule()
)

// Acker is a function we call on keys that have been processed.
type Acker func(keys ...[]byte) error

// Lazy represents an interface for lazily indexing values that have been written to DackBox.
type Lazy interface {
	Mark([]byte, proto.Message)

	Stop()
}

// NewLazy returns a new instance of a lazy indexer that reads in the values to index from the toIndex queue, indexes
// them with the given indexer, then acks indexed values with the given acker.
func NewLazy(toIndex queue.WaitableQueue, indexer Indexer, acker Acker) Lazy {
	ret := &lazyImpl{
		indexer:    indexer,
		acker:      acker,
		toIndex:    toIndex,
		toAck:      queue.NewWaitableQueue(queue.NewQueue()),
		stopSignal: concurrency.NewSignal(),
	}
	go ret.runIndexing()
	go ret.runAcking()
	return ret
}

type lazyImpl struct {
	indexer Indexer
	acker   Acker
	toIndex queue.WaitableQueue
	toAck   queue.WaitableQueue

	stopSignal concurrency.Signal
}

func (li *lazyImpl) Mark(key []byte, value proto.Message) {
	li.toIndex.Push(key, value)
}

func (li *lazyImpl) Stop() {
	li.stopSignal.Signal()
}

// No need for control logic since we always want this running with an instance of DackBox that uses lazy indexing.
func (li *lazyImpl) runIndexing() {
	for {
		select {
		case <-li.stopSignal.Done():
			return
		case <-li.toIndex.NotEmpty().Done():
		}
		// wait for queue to have contents

		key, value := li.toIndex.Pop()
		if key == nil {
			continue
		}
		if value == nil {
			err := li.indexer.Delete(key)
			if err != nil {
				log.Errorf("unable to remove value from index: %s", string(key))
			}
		} else {
			err := li.indexer.Index(key, value)
			if err != nil {
				log.Errorf("unable to add key and value to index: %s, %s", string(key), proto.MarshalTextString(value.(proto.Message)))
			}
		}
		li.toAck.Push(key, nil)
	}
}

func (li *lazyImpl) runAcking() {
	ticker := time.NewTicker(ackInterval)
	var keysToAck [][]byte
	for {
		select {
		case <-li.stopSignal.Done():
			return
		// Don't wait more than the interval to ack.
		case <-ticker.C:
			li.ackKeys(keysToAck)
			keysToAck = nil
		// Don't ack more than the max at a time.
		case <-li.toAck.NotEmpty().Done():
			key, _ := li.toAck.Pop()
			keysToAck = append(keysToAck, key)
			if len(keysToAck) == maxToAck {
				li.ackKeys(keysToAck)
				keysToAck = nil
			}
		}
	}
}

func (li *lazyImpl) ackKeys(keysToAck [][]byte) {
	if len(keysToAck) == 0 {
		return
	}
	err := li.acker(keysToAck...)
	if err != nil {
		log.Errorf("unable to ack keys: %s", printableKeys(keysToAck))
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
