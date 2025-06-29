package cache

import (
	"time"
)

func NewCache[K comparable, V any](
	cleanupInterval, cacheDuration time.Duration, actions *Actions[K, V],
) *Cache[K, V] {
	cache := Cache[K, V]{
		actions:       actions,
		cacheDuration: cacheDuration,
		cleanupChan:   make(chan bool),
		cache:         make(map[K]*CacheItem[V]),
		cleanupTicker: time.NewTicker(cleanupInterval),
	}

	if actions == nil {
		cache.actions = &Actions[K, V]{}
	}

	go cache.cleanup()
	return &cache
}

func (c *Cache[K, V]) cleanup() {
	for {
		select {
		case <-c.cleanupTicker.C:
			if len(c.cache) == 0 {
				continue
			}

			c.mutex.Lock()
			for k, v := range c.cache {
				if time.Since(v.lastAccessed) > c.cacheDuration {
					if c.actions.OnCleanup != nil {
						c.actions.OnCleanup(k, v.value)
					}
					delete(c.cache, k)
				}
			}
			c.mutex.Unlock()
		case <-c.cleanupChan:
			return
		}
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mutex.RLock()
	item, ok := c.cache[key]
	c.mutex.RUnlock()

	if !ok {
		var i V
		return i, false
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	item.lastAccessed = time.Now()

	return item.value, true
}

func (c *Cache[K, V]) Set(key K, value V) V {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, ok := c.cache[key]; ok {
		if c.actions.OnCleanup != nil {
			c.actions.OnCleanup(key, item.value)
		}
		item.value = value
		item.lastAccessed = time.Now()
	} else {
		c.cache[key] = &CacheItem[V]{
			value:        value,
			lastAccessed: time.Now(),
		}
	}

	return value
}

func (c *Cache[K, V]) Delete(key K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.cache, key)
}

func (c *Cache[K, V]) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cleanupTicker.Stop()
	close(c.cleanupChan)

	for key, item := range c.cache {
		if c.actions.OnCleanup != nil {
			c.actions.OnCleanup(key, item.value)
		}
		delete(c.cache, key)
	}
}
