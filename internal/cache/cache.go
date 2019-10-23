package cache

import (
	"sync"
)

type item struct {
	Elem interface{} `json:"elem"`
}

type Cache struct {
	mu    sync.RWMutex
	items []item
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

	c.items = append(c.items, item{
		Elem: x,
	})
}

func (c *Cache) ItemByIndex(index int) (interface{}, bool) {
	c.mu.RLock()

	defer func() {
		c.mu.RUnlock()
	}()

	if len(c.items) <= index {
		return nil, false
	}

	return c.items[index], true
}

func (c *Cache) Size() int {
	c.mu.RLock()

	defer func() {
		c.mu.RUnlock()
	}()

	return len(c.items)
}
