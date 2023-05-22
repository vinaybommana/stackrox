package lru

/*
source: https://github.com/hashicorp/golang-lru/pull/116/files#diff-4c2f540f9d9bdb85dc86612f9619b085c264052b34d125bdad879677423b4c13
This file is a temporary copy of the above referenced source. It is meant to disappear once the pending pull request
to add expiring cache is merged ( https://github.com/hashicorp/golang-lru/pull/116 ).

The lruInterface was copied and made private. The methods were made private as well.
The expirableLRU object is supposed to be package internal, so its methods were made private.
*/

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"sync"
	"testing"
	"time"
)

func BenchmarkExpirableLRU_Rand_NoExpire(b *testing.B) {
	l := newExpirableLRU[int64, int64](8192, nil, nil, 0)

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		trace[i] = getRand(b) % 32768
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < 2*b.N; i++ {
		if i%2 == 0 {
			l.add(trace[i], trace[i])
		} else {
			if _, ok := l.get(trace[i]); ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(hit+miss))
}

func BenchmarkExpirableLRU_Freq_NoExpire(b *testing.B) {
	l := newExpirableLRU[int64, int64](8192, nil, nil, 0)

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		if i%2 == 0 {
			trace[i] = getRand(b) % 16384
		} else {
			trace[i] = getRand(b) % 32768
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.add(trace[i], trace[i])
	}
	var hit, miss int
	for i := 0; i < b.N; i++ {
		if _, ok := l.get(trace[i]); ok {
			hit++
		} else {
			miss++
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(hit+miss))
}

func BenchmarkExpirableLRU_Rand_WithExpire(b *testing.B) {
	l := newExpirableLRU[int64, int64](8192, nil, nil, time.Millisecond*10)
	defer l.close()

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		trace[i] = getRand(b) % 32768
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < 2*b.N; i++ {
		if i%2 == 0 {
			l.add(trace[i], trace[i])
		} else {
			if _, ok := l.get(trace[i]); ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(hit+miss))
}

func BenchmarkExpirableLRU_Freq_WithExpire(b *testing.B) {
	l := newExpirableLRU[int64, int64](8192, nil, nil, time.Millisecond*10)
	defer l.close()

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		if i%2 == 0 {
			trace[i] = getRand(b) % 16384
		} else {
			trace[i] = getRand(b) % 32768
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.add(trace[i], trace[i])
	}
	var hit, miss int
	for i := 0; i < b.N; i++ {
		if _, ok := l.get(trace[i]); ok {
			hit++
		} else {
			miss++
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(hit+miss))
}

func TestExpirableLRUInterface(t *testing.T) {
	var _ lruCache[int, int] = &expirableLRU[int, int]{}
}

func TestExpirableLRUNoPurge(t *testing.T) {
	lc := newExpirableLRU[string, string](10, nil, nil, 0)

	lc.add("key1", "val1")
	if lc.len() != 1 {
		t.Fatalf("length differs from expected")
	}

	v, ok := lc.peek("key1")
	if v != "val1" {
		t.Fatalf("value differs from expected")
	}
	if !ok {
		t.Fatalf("should be true")
	}

	if !lc.contains("key1") {
		t.Fatalf("should contain key1")
	}
	if lc.contains("key2") {
		t.Fatalf("should not contain key2")
	}

	v, ok = lc.peek("key2")
	if v != "" {
		t.Fatalf("should be empty")
	}
	if ok {
		t.Fatalf("should be false")
	}

	if !reflect.DeepEqual(lc.keys(), []string{"key1"}) {
		t.Fatalf("value differs from expected")
	}

	if lc.resize(0) != 0 {
		t.Fatalf("evicted count differs from expected")
	}
	if lc.resize(2) != 0 {
		t.Fatalf("evicted count differs from expected")
	}
	lc.add("key2", "val2")
	if lc.resize(1) != 1 {
		t.Fatalf("evicted count differs from expected")
	}
}

func TestExpirableMultipleClose(t *testing.T) {
	lc := newExpirableLRU[string, string](10, nil, nil, 0)
	lc.close()
	// should not panic
	lc.close()
}

func TestExpirableLRUWithPurge(t *testing.T) {
	var evicted []string
	onEvicted := func(key string, value string) { evicted = append(evicted, key, value) }
	lc := newExpirableLRU(10, onEvicted, onEvicted, 150*time.Millisecond)
	defer lc.close()

	k, v, ok := lc.getOldest()
	if k != "" {
		t.Fatalf("should be empty")
	}
	if v != "" {
		t.Fatalf("should be empty")
	}
	if ok {
		t.Fatalf("should be false")
	}

	lc.add("key1", "val1")

	time.Sleep(100 * time.Millisecond) // not enough to expire
	if lc.len() != 1 {
		t.Fatalf("length differs from expected")
	}

	v, ok = lc.get("key1")
	if v != "val1" {
		t.Fatalf("value differs from expected")
	}
	if !ok {
		t.Fatalf("should be true")
	}

	time.Sleep(200 * time.Millisecond) // expire
	v, ok = lc.get("key1")
	if ok {
		t.Fatalf("should be false")
	}
	if v != "" {
		t.Fatalf("should be nil")
	}

	if lc.len() != 0 {
		t.Fatalf("length differs from expected")
	}
	if !reflect.DeepEqual(evicted, []string{"key1", "val1"}) {
		t.Fatalf("value differs from expected")
	}

	// add new entry
	lc.add("key2", "val2")
	if lc.len() != 1 {
		t.Fatalf("length differs from expected")
	}

	k, v, ok = lc.getOldest()
	if k != "key2" {
		t.Fatalf("value differs from expected")
	}
	if v != "val2" {
		t.Fatalf("value differs from expected")
	}
	if !ok {
		t.Fatalf("should be true")
	}

	// DeleteExpired, nothing deleted
	lc.deleteExpired()
	if lc.len() != 1 {
		t.Fatalf("length differs from expected")
	}
	if !reflect.DeepEqual(evicted, []string{"key1", "val1"}) {
		t.Fatalf("value differs from expected")
	}

	// Purge, cache should be clean
	lc.purge()
	if lc.len() != 0 {
		t.Fatalf("length differs from expected")
	}
	if !reflect.DeepEqual(evicted, []string{"key1", "val1", "key2", "val2"}) {
		t.Fatalf("value differs from expected")
	}
}

func TestExpirableLRUWithPurgeEnforcedBySize(t *testing.T) {
	lc := newExpirableLRU[string, string](10, nil, nil, time.Hour)
	defer lc.close()

	for i := 0; i < 100; i++ {
		i := i
		lc.add(fmt.Sprintf("key%d", i), fmt.Sprintf("val%d", i))
		v, ok := lc.get(fmt.Sprintf("key%d", i))
		if v != fmt.Sprintf("val%d", i) {
			t.Fatalf("value differs from expected")
		}
		if !ok {
			t.Fatalf("should be true")
		}
		if lc.len() > 20 {
			t.Fatalf("length should be less than 20")
		}
	}

	if lc.len() != 10 {
		t.Fatalf("length differs from expected")
	}
}

func TestExpirableLRUConcurrency(t *testing.T) {
	lc := newExpirableLRU[string, string](0, nil, nil, 0)
	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			lc.add(fmt.Sprintf("key-%d", i/10), fmt.Sprintf("val-%d", i/10))
			wg.Done()
		}(i)
	}
	wg.Wait()
	if lc.len() != 100 {
		t.Fatalf("length differs from expected")
	}
}

func TestExpirableLRUInvalidateAndEvict(t *testing.T) {
	var evicted int
	onEvict := func(_, _ string) { evicted++ }
	lc := newExpirableLRU(-1, onEvict, onEvict, 0)

	lc.add("key1", "val1")
	lc.add("key2", "val2")

	val, ok := lc.get("key1")
	if !ok {
		t.Fatalf("should be true")
	}
	if val != "val1" {
		t.Fatalf("value differs from expected")
	}
	if evicted != 0 {
		t.Fatalf("value differs from expected")
	}

	lc.remove("key1")
	if evicted != 1 {
		t.Fatalf("value differs from expected")
	}
	val, ok = lc.get("key1")
	if val != "" {
		t.Fatalf("should be empty")
	}
	if ok {
		t.Fatalf("should be false")
	}
}

func TestLoadingExpired(t *testing.T) {
	lc := newExpirableLRU[string, string](0, nil, nil, time.Millisecond*5)
	defer lc.close()

	lc.add("key1", "val1")
	if lc.len() != 1 {
		t.Fatalf("length differs from expected")
	}

	v, ok := lc.peek("key1")
	if v != "val1" {
		t.Fatalf("value differs from expected")
	}
	if !ok {
		t.Fatalf("should be true")
	}

	v, ok = lc.get("key1")
	if v != "val1" {
		t.Fatalf("value differs from expected")
	}
	if !ok {
		t.Fatalf("should be true")
	}

	time.Sleep(time.Millisecond * 10) // wait for entry to expire
	if lc.len() != 0 {
		t.Fatalf("length differs from expected")
	}

	v, ok = lc.peek("key1")
	if v != "" {
		t.Fatalf("should be empty")
	}
	if ok {
		t.Fatalf("should be false")
	}

	v, ok = lc.get("key1")
	if v != "" {
		t.Fatalf("should be empty")
	}
	if ok {
		t.Fatalf("should be false")
	}
}

func TestExpirableLRURemoveOldest(t *testing.T) {
	lc := newExpirableLRU[string, string](2, nil, nil, 0)

	k, v, ok := lc.removeOldest()
	if k != "" {
		t.Fatalf("should be empty")
	}
	if v != "" {
		t.Fatalf("should be empty")
	}
	if ok {
		t.Fatalf("should be false")
	}

	ok = lc.remove("non_existent")
	if ok {
		t.Fatalf("should be false")
	}

	lc.add("key1", "val1")
	if lc.len() != 1 {
		t.Fatalf("length differs from expected")
	}

	v, ok = lc.get("key1")
	if !ok {
		t.Fatalf("should be true")
	}
	if v != "val1" {
		t.Fatalf("value differs from expected")
	}

	if !reflect.DeepEqual(lc.keys(), []string{"key1"}) {
		t.Fatalf("value differs from expected")
	}
	if lc.len() != 1 {
		t.Fatalf("length differs from expected")
	}

	lc.add("key2", "val2")
	if !reflect.DeepEqual(lc.keys(), []string{"key1", "key2"}) {
		t.Fatalf("value differs from expected")
	}
	if lc.len() != 2 {
		t.Fatalf("length differs from expected")
	}

	k, v, ok = lc.removeOldest()
	if k != "key1" {
		t.Fatalf("value differs from expected")
	}
	if v != "val1" {
		t.Fatalf("value differs from expected")
	}
	if !ok {
		t.Fatalf("should be true")
	}

	if !reflect.DeepEqual(lc.Kkeys(), []string{"key2"}) {
		t.Fatalf("value differs from expected")
	}
	if lc.len() != 1 {
		t.Fatalf("length differs from expected")
	}
}

func ExampleExpirableLRU() {
	// make cache with 10ms TTL and 5 max keys
	cache := newExpirableLRU[string, string](5, nil, nil, time.Millisecond*10)
	// expirable cache need to be closed after used
	defer cache.close()

	// set value under key1.
	cache.add("key1", "val1")

	// get value under key1
	r, ok := cache.get("key1")

	// check for OK value
	if ok {
		fmt.Printf("value before expiration is found: %v, value: %q\n", ok, r)
	}

	// wait for cache to expire
	time.Sleep(time.Millisecond * 12)

	// get value under key1 after key expiration
	r, ok = cache.get("key1")
	fmt.Printf("value after expiration is found: %v, value: %q\n", ok, r)

	// set value under key2, would evict old entry because it is already expired.
	cache.add("key2", "val2")

	fmt.Printf("Cache len: %d\n", cache.len())
	// Output:
	// value before expiration is found: true, value: "val1"
	// value after expiration is found: false, value: ""
	// Cache len: 1
}

func getRand(tb testing.TB) int64 {
	out, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		tb.Fatal(err)
	}
	return out.Int64()
}
