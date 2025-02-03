package pokecache

import (
	"sync"
	"time"
)

type cacheEnrty struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	caches map[string]cacheEnrty
	mu     sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	caches := make(map[string]cacheEnrty)

	cache := &Cache{
		caches: caches,
		mu:     sync.Mutex{},
	}

	cache.reapLoop(interval)

	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.caches[key] = cacheEnrty{createdAt: time.Now(), val: val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if value, ok := c.caches[key]; ok {
		return value.val, true
	}

	return nil, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			c.mu.Lock()

			for key, value := range c.caches {
				if time.Since(value.createdAt) > interval {
					delete(c.caches, key)
				}
			}

			c.mu.Unlock()
		}
	}()
}
