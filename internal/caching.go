package internal

import (
	"sync"
	"time"
)

type Cache struct {
	entry    map[string]cacheEntry
	interval time.Duration
	mutex    *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c *Cache) Add(key string, val []byte) {
	NewEntry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.mutex.Lock()
	c.entry[key] = NewEntry
	c.mutex.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.entry[key]
	if !exists {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) ReapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.Reap(time.Now().UTC(), interval)

	}
}

func (c *Cache) Reap(now time.Time, last time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key, val := range c.entry {
		if val.createdAt.Before(now.Add(-last)) {
			delete(c.entry, key)
		}
	}
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entry:    make(map[string]cacheEntry),
		interval: interval,
		mutex:    &sync.Mutex{},
	}

	go c.ReapLoop(interval)

	return c
}
