package Cache

import (
	"Cache/lru"
	"sync"
)

type cache struct {
	mu      sync.Mutex
	lru     *lru.Lru
	maxSize int
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.maxSize, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (ByteView, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return ByteView{}, false
	}
	res, ok := c.lru.Get(key)
	if !ok {
		return ByteView{}, false
	}
	return res.(ByteView), ok

}
