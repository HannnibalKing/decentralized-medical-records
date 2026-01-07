package revocation

import (
	"sync"
	"time"
)

// Cache provides a tiny in-memory revocation cache with TTL.
type Cache struct {
	mu    sync.RWMutex
	items map[string]entry
	ttl   time.Duration
}

type entry struct {
	revoked bool
	exp     time.Time
}

// NewCache builds a cache with the given TTL.
func NewCache(ttl time.Duration) *Cache {
	return &Cache{items: make(map[string]entry), ttl: ttl}
}

// Set records a revocation state.
func (c *Cache) Set(handle string, revoked bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[handle] = entry{revoked: revoked, exp: time.Now().Add(c.ttl)}
}

// Get returns (revoked, ok).
func (c *Cache) Get(handle string) (bool, bool) {
	c.mu.RLock()
	e, ok := c.items[handle]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.exp) {
		if ok {
			c.mu.Lock()
			delete(c.items, handle)
			c.mu.Unlock()
		}
		return false, false
	}
	return e.revoked, true
}
