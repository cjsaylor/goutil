// Package lru is a package that implements a "least recently used" data structure.
package lru

import (
	"container/list"
	"sync"
)

type entry struct {
	key   interface{}
	value interface{}
}

// EvictionCallback is a method you can specify to receive evicted values from the LRU cache.
type EvictionCallback func(key, value interface{})

// Cache is a key-value store with a fixed length. The oldest entry will be evicted when the newest entry
// is added at the capacity limit.
type Cache struct {
	queue      *list.List
	lookup     map[interface{}]*list.Element
	capacity   int
	mutex      *sync.Mutex
	onEviction EvictionCallback
}

// NewCache creates an instance of an LRU cache with fixed capacity.
func NewCache(capacity int, onEviction EvictionCallback) *Cache {
	cache := Cache{
		queue:      list.New(),
		lookup:     make(map[interface{}]*list.Element, capacity),
		capacity:   capacity,
		mutex:      &sync.Mutex{},
		onEviction: onEviction,
	}
	return &cache
}

// Set a key/value into the LRU cache.
// This will evict the oldest entry if at the capacity limit.
func (c *Cache) Set(key, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if item, ok := c.lookup[key]; ok {
		c.queue.MoveToFront(item)
		item.Value = &entry{
			key:   key,
			value: value,
		}
	}
	item := c.queue.PushFront(&entry{
		key:   key,
		value: value,
	})
	c.lookup[key] = item
	if c.queue.Len() > c.capacity {
		c.RemoveOldest()
	}
}

// Get will retrieve a value by key.
// This will bump the entry as it was "recently" used.
func (c *Cache) Get(key interface{}) (interface{}, bool) {
	if item, ok := c.lookup[key]; ok {
		c.queue.MoveToFront(item)
		return item.Value, true
	}
	return nil, false
}

// Remove an entry from the LRU cache
func (c *Cache) Remove(key interface{}) (interface{}, bool) {
	if item, ok := c.lookup[key]; ok {
		c.queue.Remove(item)
		delete(c.lookup, key)
		return item.Value.(*entry).value, true
	}
	return nil, false
}

// RemoveOldest will remove the oldest entry from the LRU cache.
func (c *Cache) RemoveOldest() (interface{}, bool) {
	if c.queue.Len() == 0 {
		return nil, false
	}
	tail := c.queue.Back()
	c.queue.Remove(tail)
	delete(c.lookup, tail.Value.(*entry).key)
	c.onEviction(tail.Value.(*entry).key, tail.Value.(*entry).value)
	return tail.Value.(*entry).value, true
}
