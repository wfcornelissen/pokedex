package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entry    map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	return &Cache{
		Entry:    make(map[string]cacheEntry),
		interval: interval,
	}
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Entry[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.Entry[key]
	if exists {
		return entry.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop() {

}
