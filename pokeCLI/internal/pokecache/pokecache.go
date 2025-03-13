package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	EntryMap map[string]cacheEntry
	Mu       sync.Mutex
	Interval time.Duration
}

type CacheInterface interface {
	NewCache() Cache
	Add()
	Get() ([]byte, error)
	reapLoop()
}

type cacheEntry struct {
	CreatedAt time.Time // Represents when the entry was created.
	Val       []byte    // Represents the raw data we're caching.
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		EntryMap: make(map[string]cacheEntry),
		Mu:       sync.Mutex{},
		Interval: interval,
	}
	go cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	c.EntryMap[key] = cacheEntry{
		CreatedAt: time.Now(),
		Val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	entry, ok := c.EntryMap[key]
	return entry.Val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.Interval)
	defer ticker.Stop()

	c.Mu.Lock()
	for key, val := range c.EntryMap {
		if time.Since(val.CreatedAt) > c.Interval {
			delete(c.EntryMap, key)
		}
	}
	c.Mu.Unlock()
}
