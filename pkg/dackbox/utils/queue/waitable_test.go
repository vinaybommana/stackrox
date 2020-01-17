package queue

import (
	"testing"
	"time"

	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stretchr/testify/assert"
)

func TestWaitableQueue(t *testing.T) {
	q := NewWaitableQueue(NewQueue())

	input := [][]byte{
		[]byte("id1"),
		[]byte("id2"),
		[]byte("id3"),
		[]byte("id4"),
	}
	q.Push(input[0], nil)
	q.Push(input[1], nil)
	q.Push(input[2], nil)
	q.Push(input[3], nil)

	var results [][]byte
	for q.Length() > 0 {
		k, _ := q.Pop()
		results = append(results, k)
	}
	assert.Equal(t, input, results)

	q.Push(input[0], nil)
	q.Push(input[1], nil)
	q.Push(input[2], nil)
	q.Push(input[3], nil)

	results = [][]byte{}
	for q.Length() > 0 {
		k, _ := q.Pop()
		results = append(results, k)
	}

	assert.Equal(t, input, results)
}

func TestWaitableQueueConcurrent(t *testing.T) {
	q := NewWaitableQueue(NewQueue())

	input := [][]byte{
		[]byte("id1"),
		[]byte("id2"),
		[]byte("id3"),
		[]byte("id4"),
	}

	assertable := concurrency.NewSignal()
	var results [][]byte
	go func() {
		for len(results) < 4 {
			<-q.NotEmpty().Done()
			k, _ := q.Pop()
			results = append(results, k)
		}
		assertable.Signal()
	}()
	q.Push(input[0], nil)
	q.Push(input[1], nil)
	time.Sleep(5 * time.Millisecond)
	q.Push(input[2], nil)
	q.Push(input[3], nil)
	select {
	case <-time.NewTimer(time.Second).C:
		assert.Fail(t, "assertable never returned")
	case <-assertable.Done():
	}

	assert.Equal(t, input, results)
}
