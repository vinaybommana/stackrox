package throttle

import (
	"testing"
	"time"

	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stretchr/testify/assert"
)

func TestThrottlesFastCalls(t *testing.T) {
	throttler := NewDropThrottle(10 * time.Millisecond)

	// Run count should be two, one for the first, and one for the end of the window since more were called.
	var ran string
	var ranMutex sync.Mutex
	throttler.Run(func() {
		concurrency.WithLock(&ranMutex, func() {
			ran = ran + "1"
		})
	})

	throttler.Run(func() {
		concurrency.WithLock(&ranMutex, func() {
			ran = ran + "2"
		})
	})

	throttler.Run(func() {
		concurrency.WithLock(&ranMutex, func() {
			ran = ran + "3"
		})
	})

	throttler.Run(func() {
		concurrency.WithLock(&ranMutex, func() {
			ran = ran + "4"
		})
	})

	time.Sleep(20 * time.Millisecond)

	ranMutex.Lock()
	defer ranMutex.Unlock()
	assert.Equal(t, "12", ran)
}

func TestThrottlesSlowCalls(t *testing.T) {
	throttler := NewDropThrottle(10 * time.Millisecond)

	// Run count should be two, one for the first, and one for the end of the window since more were called.
	var ran string
	var ranMutex sync.Mutex
	throttler.Run(func() {
		time.Sleep(20 * time.Millisecond)
		concurrency.WithLock(&ranMutex, func() {
			ran = ran + "1"
		})
	})

	throttler.Run(func() {
		time.Sleep(20 * time.Millisecond)
		concurrency.WithLock(&ranMutex, func() {
			ran = ran + "2"
		})
	})

	throttler.Run(func() {
		time.Sleep(20 * time.Millisecond)
		concurrency.WithLock(&ranMutex, func() {
			ran = ran + "3"
		})
	})

	throttler.Run(func() {
		time.Sleep(20 * time.Millisecond)
		concurrency.WithLock(&ranMutex, func() {
			ran = ran + "4"
		})
	})

	// Both should complete by 30 millis.
	// the first will Run in 20, and the last will Run in 20 after a 10 millisecond wait window)
	time.Sleep(40 * time.Millisecond)
	ranMutex.Lock()
	defer ranMutex.Unlock()
	assert.Equal(t, "12", ran)
}
