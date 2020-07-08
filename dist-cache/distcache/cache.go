package distcache

import (
	"proverbs.top/distcache/lru"
	"sync"
)

type Cache struct {
	mu sync.Mutex
	lru_cache *lru.LRUCache
	cacheBytes int64
}

func (c *Cache) add(key string, val ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru_cache == nil {
		c.lru_cache = lru.New(c.cacheBytes, nil)
	}
	c.lru_cache.Add(key, val)
}

func (c *Cache) get(key string) (ByteView, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru_cache == nil {
		return ByteView{}, false
	}
	if v, ok := c.lru_cache.Get(key); ok {
		return v.(ByteView), true
	}
	return ByteView{}, false
}
