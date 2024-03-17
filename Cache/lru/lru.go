// Package lru
package lru

import "container/list"

type Value interface {
	Len() int
}
type entry struct {
	key   string
	value Value
}

type Lru struct {
	maxBytes int
	nbytes   int
	cache    map[string]*list.Element
	ll       *list.List
	// optional, you can implement it in different ways.
	onEvicted func(key string, value Value)
}

func New(maxBytes int, onEvicted func(key string, value Value)) *Lru {
	return &Lru{
		maxBytes:  maxBytes,
		nbytes:    0,
		cache:     make(map[string]*list.Element),
		ll:        list.New(),
		onEvicted: onEvicted,
	}
}

func (l *Lru) Get(key string) (Value, bool) {
	e, ok := l.cache[key]
	if ok {
		l.ll.MoveToFront(e)
		return e.Value.(*entry).value, ok
	}
	return nil, ok
}

func (l *Lru) RemoveOldest() {
	e := l.ll.Back()
	if e != nil {
		l.ll.Remove(e)
		kv := e.Value.(*entry)
		delete(l.cache, kv.key)
		l.nbytes -= kv.value.Len() + len(kv.key)
		if l.onEvicted != nil {
			l.onEvicted(kv.key, kv.value)
		}
	}
}

func (l *Lru) Add(key string, value Value) {
	if e, ok := l.cache[key]; ok {
		l.ll.MoveToFront(e)
		kv := e.Value.(*entry)
		e.Value.(*entry).value = value
		l.nbytes += value.Len() - kv.value.Len()
	} else {
		e := l.ll.PushFront(&entry{key, value})
		l.cache[key] = e
		l.nbytes += len(key) + value.Len()
	}
	if l.maxBytes != 0 && l.maxBytes < l.nbytes {
		l.RemoveOldest()
	}
}

func (l *Lru) Len() int {
	return l.ll.Len()
}
