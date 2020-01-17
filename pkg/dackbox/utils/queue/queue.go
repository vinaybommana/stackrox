package queue

import (
	"github.com/gogo/protobuf/proto"
)

// Queue holds a deduplicated set of keys and values.
// If a new value is pushed for a key currently in the queue, the key will maintain its position, but the value will be
// updated.
// This is NOT thread-safe.
type Queue interface {
	Push(key []byte, value proto.Message)
	Pop() ([]byte, proto.Message)
	Contains(key []byte) bool
	Length() int
}

// NewQueue returns a new instance of a Queue
func NewQueue() Queue {
	return &queueImpl{
		kvPairs: make(map[string]proto.Message),
	}
}

type queueImpl struct {
	kvPairs map[string]proto.Message
	front   *keyNode
	back    *keyNode
	length  int
}

func (q *queueImpl) Push(key []byte, value proto.Message) {
	keyString := string(key)
	_, exists := q.kvPairs[keyString]
	q.kvPairs[keyString] = value
	if exists { // if the key is already in the queue, just update the value.
		return
	}
	q.pushKey(key)
}

func (q *queueImpl) Pop() ([]byte, proto.Message) {
	key := q.popKey()
	if key == nil {
		return nil, nil
	}
	keyString := string(key)
	value := q.kvPairs[keyString]
	delete(q.kvPairs, keyString)
	return key, value
}

func (q *queueImpl) Contains(key []byte) bool {
	_, exists := q.kvPairs[string(key)]
	return exists
}

func (q *queueImpl) Length() int {
	return q.length
}

type keyNode struct {
	key  []byte
	next *keyNode
	prev *keyNode
}

func (q *queueImpl) pushKey(key []byte) {
	newNode := &keyNode{
		key:  key,
		next: q.back,
	}
	// Set new node as back
	if q.back != nil {
		q.back.prev = newNode
	}
	q.back = newNode
	// If queue was empty, new node is now front as well.
	if q.front == nil {
		q.front = newNode
	}
	q.length++
}

func (q *queueImpl) popKey() []byte {
	if q.front == nil {
		return nil
	}
	// Get key from front.
	ret := q.front.key
	// set front to it's previous value.
	q.front = q.front.prev
	// If front exists, reset it's next value to null.
	if q.front != nil {
		q.front.next = nil
	} else {
		q.back = nil
	}
	q.length--
	return ret
}
