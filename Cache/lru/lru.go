package cache

import "container/list"

type Cache struct {
	// max cache
	maxBytes int
	// used
	nbytes int
	// hash map , preserve the position of each element.
	cache map[string]*list.Element
	//
	ll *list.List
	// callback function.
	OnEvicted func(key string, value Value)
}
type entry struct {
	key   string
	value Value
}
type Value interface {
	Len() int
}

// New init a Cache
func New(maxBytes int, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get get element and move element to front of list.
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ent, ok := c.cache[key]; ok {
		// move to front
		c.ll.MoveToFront(ent)
		return ent.Value.(*entry).value, true
	}
	return
}

// RemoveOldest remove the oldest of unused element.
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int(kv.value.Len()) + int(len(kv.key))
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add add new element to front of list.
func (c *Cache) Add(key string, value Value) {
	if ent, ok := c.cache[key]; ok {
		// update
		c.ll.MoveToFront(ent)
		kv := ent.Value.(*entry)
		c.nbytes += int(value.Len()) - int(kv.value.Len())
		kv.value = value
		return
	}
	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele
	c.nbytes += int(value.Len()) + int(len(key))
	if c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}

}

// Len cache length
func (c *Cache) Len() int {
	return c.ll.Len()
}
