package lru

import (
	hashicorpLRU "github.com/hashicorp/golang-lru/v2"
)

type Cache[K comparable, V any] struct {
	*hashicorpLRU.Cache[K, V]
}

func New[K comparable, V any](size int) (*Cache[K, V], error) {
	cache, err := hashicorpLRU.New[K, V](size)
	return &Cache[K, V]{Cache: cache}, err
}

func NewWithEvict[K comparable, V any](size int, onEvicted func(key K, value V)) (*Cache[K, V], error) {
	cache, err := hashicorpLRU.NewWithEvict[K, V](size, onEvicted)
	return &Cache[K, V]{Cache: cache}, err
}
