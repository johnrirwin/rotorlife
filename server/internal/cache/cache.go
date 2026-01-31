package cache

import (
	"sync"
	"time"
)

// MemoryCache is an in-memory cache implementation with TTL support
type MemoryCache struct {
	mu     sync.RWMutex
	items  map[string]entry
	ttl    time.Duration
	stopCh chan struct{}
}

type entry struct {
	value     interface{}
	expiresAt time.Time
}

// NewMemory creates a new in-memory cache with the specified TTL
func NewMemory(ttl time.Duration) *MemoryCache {
	c := &MemoryCache{
		items:  make(map[string]entry),
		ttl:    ttl,
		stopCh: make(chan struct{}),
	}
	go c.cleanup()
	return c
}

// New creates a new in-memory cache (alias for NewMemory for backwards compatibility)
func New(ttl time.Duration) *MemoryCache {
	return NewMemory(ttl)
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	e, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e.value, true
}

func (c *MemoryCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = entry{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *MemoryCache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = entry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *MemoryCache) Invalidate(key string) {
	c.Delete(key)
}

func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]entry)
}

func (c *MemoryCache) Stop() {
	close(c.stopCh)
}

func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.removeExpired()
		case <-c.stopCh:
			return
		}
	}
}

func (c *MemoryCache) removeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, e := range c.items {
		if now.After(e.expiresAt) {
			delete(c.items, key)
		}
	}
}

// Ensure MemoryCache implements Cache interface
var _ Cache = (*MemoryCache)(nil)
