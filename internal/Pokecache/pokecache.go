package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	CacheEntry map[string]cacheEntry
	mu         sync.RWMutex
}
type cacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		CacheEntry: make(map[string]cacheEntry),
		mu:         sync.RWMutex{},
	}

	cache.reapLoop(interval)

	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.CacheEntry[key] = cacheEntry{time.Now(), val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.CacheEntry[key]; ok {
		return entry.Val, true
	}

	return nil, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Tick at", time.Now())
			case <-done:
				c.CacheEntry = make(map[string]cacheEntry)
			}
		}
	}()

	time.Sleep(interval)
	done <- true
}
