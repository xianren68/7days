package singleflight

import "sync"

// Group manage request.
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}
type call struct {
	wg  sync.WaitGroup
	val any
	err error
}

// Do ensure only one request for the same key.
func (g *Group) Do(key string, fn func() (any, error)) (any, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// exist request.
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		// wait request response.
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
