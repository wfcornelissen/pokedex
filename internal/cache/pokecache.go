package pokecache

import (
	"fmt"
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
	c := &Cache{
		Entry:    make(map[string]cacheEntry),
		interval: interval * time.Second,
	}
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key *string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Entry[*key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key *string) ([]byte, bool) {
	if key == nil {
		return nil, false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.Entry[*key]
	if exists {
		return entry.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop() {
	fmt.Println("reapLoop called")
	for {
		c.mu.Lock()
		for url, entry := range c.Entry {
			if time.Since(entry.createdAt) > c.interval {
				delete(c.Entry, url)
			}
		}
		c.mu.Unlock()
		time.Sleep(c.interval)
	}
}
