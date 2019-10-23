package cache

import (
	"sync"
)

// todo: add expiration

type Item struct {
	Object interface{}
}

type Cache struct {
	mu    sync.RWMutex
	items []Item
}

func New() *Cache {
	cache := Cache{}
	return &cache
}

func (c *Cache) AddItem(x interface{}) {
	c.mu.Lock()

	defer func() {
		c.mu.Unlock()
	}()

	c.items = append(c.items, Item{
		Object: x,
	})
}

func (c *Cache) ItemByIndex(index int) (interface{}, bool) {
	c.mu.RLock()

	defer func() {
		c.mu.Unlock()
	}()

	if len(c.items) <= index {
		return nil, false
	} else {
		return c.items[index], true
	}
}

func (c *Cache) Size() int {
	c.mu.RLock()

	defer func() {
		c.mu.RUnlock()
	}()

	return len(c.items)
}