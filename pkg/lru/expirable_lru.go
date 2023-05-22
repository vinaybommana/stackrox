package lru

/*
source: https://github.com/hashicorp/golang-lru/pull/116/files#diff-be4b9429f63e821595d57501426a7f91505ab8f2df186c5f5e3b0af74e1a5dfc
This file is a temporary copy of the above referenced source. It is meant to disappear once the pending pull request
to add expiring cache is merged ( https://github.com/hashicorp/golang-lru/pull/116 ).

The methods of this LRU cache internals implementation were made private as the expirableLRU should be wrapped in
a public Cache which will then expose other properties.
*/

import (
	"time"

	"github.com/stackrox/rox/pkg/sync"
)

// EvictCallback is used to get a callback when a cache entry is evicted
type EvictCallback[K comparable, V any] func(key K, value V)

// expirableLRU implements a thread-safe LRU with expirable entries.
//
// Entries can be cleaned up from cache with up to 1% of ttl unused.
// It happens because cleanup mechanism puts them 99 cleanup buckets away
// from the current moment, and then cleans them up 99% of ttl later instead of 100%.
type expirableLRU[K comparable, V any] struct {
	size      int
	evictList *lruList[K, V]
	items     map[K]*entry[K, V]
	onEvict   EvictCallback[K, V]
	onExpire  EvictCallback[K, V]

	// expirable options
	mu   sync.Mutex
	ttl  time.Duration
	done chan struct{}

	// buckets for expiration
	buckets []bucket[K, V]
	// uint8 because it's number between 0 and numBuckets
	nextCleanupBucket uint8
}

// bucket is a container for holding entries to be expired
type bucket[K comparable, V any] struct {
	entries map[K]*entry[K, V]
}

// noEvictionTTL - very long ttl to prevent eviction
const noEvictionTTL = time.Hour * 24 * 365 * 10

// because of uint8 usage for nextCleanupBucket, should not exceed 256.
const numBuckets = 100

// newExpirableLRU returns a new thread-safe cache with expirable entries.
//
// Size parameter set to 0 makes cache of unlimited size, e.gg turns LRU mechanism off.
//
// Providing 0 TTL turns expiring off.
//
// Delete expired entries every 1/100th of ttl value.
func newExpirableLRU[K comparable, V any](size int, onEvict EvictCallback[K, V], onExpire EvictCallback[K, V], ttl time.Duration) *expirableLRU[K, V] {
	if size < 0 {
		size = 0
	}
	if ttl <= 0 {
		ttl = noEvictionTTL
	}

	res := expirableLRU[K, V]{
		ttl:       ttl,
		size:      size,
		evictList: newList[K, V](),
		items:     make(map[K]*entry[K, V]),
		onEvict:   onEvict,
		onExpire:  onExpire,
		done:      make(chan struct{}),
	}

	// initialize the buckets
	res.buckets = make([]bucket[K, V], numBuckets)
	for i := 0; i < numBuckets; i++ {
		res.buckets[i] = bucket[K, V]{entries: make(map[K]*entry[K, V])}
	}

	// enable deleteExpired() running in separate goroutine for cache
	// with non-zero TTL
	if res.ttl != noEvictionTTL {
		go func(done <-chan struct{}) {
			ticker := time.NewTicker(res.ttl / numBuckets)
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					res.deleteExpired()
				}
			}
		}(res.done)
	}
	return &res
}

// purge clears the cache completely.
// onEvict is called for each evicted key.
func (c *expirableLRU[K, V]) purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.items {
		if c.onEvict != nil {
			c.onEvict(k, v.value)
		}
		delete(c.items, k)
	}
	for _, b := range c.buckets {
		for _, ent := range b.entries {
			delete(b.entries, ent.key)
		}
	}
	c.evictList.init()
}

// add adds a value to the cache. Returns true if an eviction occurred.
// Returns false if there was no eviction: the item was already in the cache,
// or the size was not exceeded.
func (c *expirableLRU[K, V]) add(key K, value V) (evicted bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()

	// Check for existing item
	if ent, ok := c.items[key]; ok {
		c.evictList.moveToFront(ent)
		c.removeFromBucket(ent) // remove the entry from its current bucket as expiresAt is renewed
		ent.value = value
		ent.expiresAt = now.Add(c.ttl)
		c.addToBucket(ent)
		return false
	}

	// Add new item
	ent := c.evictList.pushFrontExpirable(key, value, now.Add(c.ttl))
	c.items[key] = ent
	c.addToBucket(ent) // add the entry to the appropriate bucket and sets entry.expiredBucket

	evict := c.size > 0 && c.evictList.length() > c.size
	// Verify size not exceeded
	if evict {
		c.removeOldestNoLock()
	}
	return evict
}

// get looks up a key's value from the cache.
func (c *expirableLRU[K, V]) get(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var ent *entry[K, V]
	if ent, ok = c.items[key]; ok {
		// Expired item check
		if time.Now().After(ent.expiresAt) {
			return
		}
		c.evictList.moveToFront(ent)
		return ent.value, true
	}
	return
}

// contains checks if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (c *expirableLRU[K, V]) contains(key K) (ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok = c.items[key]
	return ok
}

// peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *expirableLRU[K, V]) peek(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var ent *entry[K, V]
	if ent, ok = c.items[key]; ok {
		// Expired item check
		if time.Now().After(ent.expiresAt) {
			return
		}
		return ent.value, true
	}
	return
}

// remove removes the provided key from the cache, returning if the
// key was contained.
func (c *expirableLRU[K, V]) remove(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
		return true
	}
	return false
}

// removeOldest removes the oldest item from the cache.
func (c *expirableLRU[K, V]) removeOldest() (key K, value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ent := c.evictList.back(); ent != nil {
		c.removeElement(ent)
		return ent.key, ent.value, true
	}
	return
}

// getOldest returns the oldest entry
func (c *expirableLRU[K, V]) getOldest() (key K, value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ent := c.evictList.back(); ent != nil {
		return ent.key, ent.value, true
	}
	return
}

// keys returns a slice of the keys in the cache, from the oldest to the newest.
func (c *expirableLRU[K, V]) keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := make([]K, 0, len(c.items))
	for ent := c.evictList.back(); ent != nil; ent = ent.prevEntry() {
		keys = append(keys, ent.key)
	}
	return keys
}

// len returns the number of items in the cache.
func (c *expirableLRU[K, V]) len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.evictList.length()
}

// resize changes the cache size. Size of 0 means unlimited.
func (c *expirableLRU[K, V]) resize(size int) (evicted int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if size <= 0 {
		c.size = 0
		return 0
	}
	diff := c.evictList.length() - size
	if diff < 0 {
		diff = 0
	}
	for i := 0; i < diff; i++ {
		c.removeOldestNoLock()
	}
	c.size = size
	return diff
}

// close destroys the cleanup goroutine. To clean up the cache, run purge() before close().
func (c *expirableLRU[K, V]) close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	select {
	case <-c.done:
		return
	default:
	}
	close(c.done)
}

// removeOldestNoLock removes the oldest item from the cache. Has to be called with lock!
func (c *expirableLRU[K, V]) removeOldestNoLock() {
	if ent := c.evictList.back(); ent != nil {
		c.removeElement(ent)
	}
}

// removeElement is used to remove a given list element from the cache. Has to be called with lock!
func (c *expirableLRU[K, V]) removeElement(e *entry[K, V]) {
	c.evictList.remove(e)
	delete(c.items, e.key)
	c.removeFromBucket(e)
	if c.onEvict != nil {
		c.onEvict(e.key, e.value)
	}
}

// deleteExpired deletes expired records. Doesn't check for entry.expiresAt as it could be
// TTL/numBuckets in the future, with numBuckets of 100 its 1% of wasted TTL.
func (c *expirableLRU[K, V]) deleteExpired() {
	c.mu.Lock()
	bucketIdx := c.nextCleanupBucket
	evictedEntries := make([]*entry[K, V], 0, len(c.buckets[bucketIdx].entries))
	for _, ent := range c.buckets[bucketIdx].entries {
		evictedEntries = append(evictedEntries, ent)
		c.evictList.remove(ent)
		delete(c.items, ent.key)
		c.removeFromBucket(ent)
	}
	c.nextCleanupBucket = (c.nextCleanupBucket + 1) % numBuckets
	c.mu.Unlock()
	for _, ent := range evictedEntries {
		c.onExpire(ent.key, ent.value)
	}
}

// addToBucket adds entry to expire bucket so that it will be cleaned up when the time comes. Has to be called with lock!
func (c *expirableLRU[K, V]) addToBucket(e *entry[K, V]) {
	bucketID := (numBuckets + c.nextCleanupBucket - 1) % numBuckets
	e.expireBucket = bucketID
	c.buckets[bucketID].entries[e.key] = e
}

// removeFromBucket removes the entry from its corresponding bucket. Has to be called with lock!
func (c *expirableLRU[K, V]) removeFromBucket(e *entry[K, V]) {
	delete(c.buckets[e.expireBucket].entries, e.key)
}
