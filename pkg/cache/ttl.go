package cache

import (
	"sync"
	"time"
)

type entry[V any] struct {
	value     V
	expiresAt time.Time
}

// TTL é uma cache em memória genérica com expiração por entrada.
// Thread-safe via sync.Map. Zero-value não é utilizável — usar New.
type TTL[K comparable, V any] struct {
	mu  sync.Map
	ttl time.Duration
}

func New[K comparable, V any](ttl time.Duration) *TTL[K, V] {
	return &TTL[K, V]{ttl: ttl}
}

func (c *TTL[K, V]) Get(key K) (V, bool) {
	raw, ok := c.mu.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	e := raw.(entry[V])
	if time.Now().After(e.expiresAt) {
		c.mu.Delete(key)
		var zero V
		return zero, false
	}
	return e.value, true
}

func (c *TTL[K, V]) Set(key K, value V) {
	c.mu.Store(key, entry[V]{value: value, expiresAt: time.Now().Add(c.ttl)})
}

func (c *TTL[K, V]) Delete(key K) {
	c.mu.Delete(key)
}
