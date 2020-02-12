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

	assertableSignal := concurrency.NewSignal()
	var results [][]byte
	go func() {
		for len(results) < 4 {
			// Wait for q to have values.
			<-q.NotEmpty().Done()

			// Pop the next value.
			k, _ := q.Pop()
			results = append(results, k)
			assertableSignal.Signal()
		}
	}()

	// Add values to the queue.
	q.Push(input[0], nil)
	q.Push(input[1], nil)

	// Wait for the q to empty.
	select {
	case <-time.After(time.Second):
		assert.Fail(t, "assertable never returned")
	case <-q.Empty().Done():
	}

	// Add more values to the queue.
	q.Push(input[2], nil)
	q.Push(input[3], nil)

	// Wait for the q to empty again.
	select {
	case <-time.After(time.Second):
		assert.Fail(t, "assertable never returned")
	case <-q.Empty().Done():
	}

	// Wait for the thread building the results to be done.
	select {
	case <-time.After(time.Second):
		assert.Fail(t, "assertable never returned")
	case <-assertableSignal.Done():
	}

	assert.Equal(t, input, results)
}
