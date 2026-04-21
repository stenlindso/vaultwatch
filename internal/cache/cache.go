// Package cache provides a simple time-based in-memory cache for Vault path listings.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached list of paths along with the time it was stored.
type Entry struct {
	Paths     []string
	CachedAt  time.Time
}

// Cache is a thread-safe in-memory store for path listings keyed by environment name.
type Cache struct {
	mu      sync.RWMutex
	store   map[string]Entry
	ttl     time.Duration
}

// New creates a new Cache with the given TTL. Entries older than ttl are
// considered stale and will not be returned by Get.
func New(ttl time.Duration) *Cache {
	return &Cache{
		store: make(map[string]Entry),
		ttl:   ttl,
	}
}

// Set stores the paths for the given environment key.
func (c *Cache) Set(env string, paths []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[env] = Entry{
		Paths:    paths,
		CachedAt: time.Now(),
	}
}

// Get returns the cached paths for env if they exist and are not stale.
// The second return value reports whether a valid entry was found.
func (c *Cache) Get(env string) ([]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.store[env]
	if !ok {
		return nil, false
	}
	if time.Since(entry.CachedAt) > c.ttl {
		return nil, false
	}
	return entry.Paths, true
}

// Invalidate removes the cached entry for the given environment.
func (c *Cache) Invalidate(env string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, env)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]Entry)
}

// Size returns the number of entries currently in the cache (including stale ones).
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.store)
}
