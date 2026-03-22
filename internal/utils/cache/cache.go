package cache

import (
	"hash/maphash"
	"sync"
)

type cache[T any] map[uint64]T

type CacheManager[T any] struct {
	mu    sync.RWMutex
	store cache[T]
	seed  maphash.Seed
}

func NewCache[T any]() *CacheManager[T] {
	return &CacheManager[T]{
		store: make(cache[T]),
		seed:  maphash.MakeSeed(),
	}
}

func (c *CacheManager[T]) Hash(strs ...string) uint64 {
	var h maphash.Hash
	h.SetSeed(c.seed)
	for _, s := range strs {
		h.WriteString(s)
		// Preventing collisions
		h.WriteByte(0)
	}
	return h.Sum64()
}

func (c *CacheManager[T]) Get(strs ...string) (T, bool) {
	key := c.Hash(strs...)
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.store[key]
	return val, ok
}

func (c *CacheManager[T]) Set(value T, strs ...string) {
	key := c.Hash(strs...)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}
