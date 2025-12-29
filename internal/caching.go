package internal

import (
	"sync"
	"time"
)

type cache struct {
	entry    map[string]cacheEntry
	interval time.Duration
	mutex    *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c *cache) Add(key string, val []byte) {
	NewEntry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.mutex.Lock()
	c.entry[key] = NewEntry
	c.mutex.Unlock()
}

func (c *cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.entry[key]
	if !exists {
		return nil, false
	}

	// Check if entry has expired
	if time.Since(entry.createdAt) > c.interval {
		delete(c.entry, key)
		return nil, false
	}

	return entry.val, true
}

func (c *cache) ReapLoop(ticker *time.Ticker) {
	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, entry := range c.entry {
			if now.Sub(entry.createdAt) > c.interval {
				delete(c.entry, key)
			}
		}
		c.mutex.Unlock()
	}
}

func NewCache(interval time.Duration) *cache {
	c := &cache{
		entry:    make(map[string]cacheEntry),
		interval: interval,
		mutex:    &sync.Mutex{},
	}
	ticker := time.NewTicker(interval)
	go c.ReapLoop(ticker)

	return c
}
