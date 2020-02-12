package queue

import (
	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/sync"
)

// WaitableQueue is a thread safe queue with an extra provided function that allows you to wait for a value to pop.
type WaitableQueue interface {
	Queue

	NotEmpty() concurrency.Waitable
	Empty() concurrency.Waitable
}

// NewWaitableQueue return a new instance of a WaitableQueue.
func NewWaitableQueue(base Queue) WaitableQueue {
	return &waitableQueueImpl{
		notEmptySig: concurrency.NewSignal(),
		base:        base,
	}
}

type waitableQueueImpl struct {
	lock        sync.Mutex
	emptySig    concurrency.Signal
	notEmptySig concurrency.Signal
	base        Queue
}

func (q *waitableQueueImpl) NotEmpty() concurrency.Waitable {
	return q.notEmptySig.WaitC()
}

func (q *waitableQueueImpl) Empty() concurrency.Waitable {
	return q.emptySig.WaitC()
}

func (q *waitableQueueImpl) Push(key []byte, value proto.Message) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.notEmptySig.Signal()
	q.emptySig.Reset()

	q.base.Push(key, value)
}

func (q *waitableQueueImpl) Pop() (key []byte, value proto.Message) {
	q.lock.Lock()
	defer q.lock.Unlock()

	key, value = q.base.Pop()
	if q.base.Length() == 0 {
		q.notEmptySig.Reset()
		q.emptySig.Signal()
	}
	return key, value
}

func (q *waitableQueueImpl) Contains(key []byte) bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.base.Contains(key)
}

func (q *waitableQueueImpl) Length() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.base.Length()
}
