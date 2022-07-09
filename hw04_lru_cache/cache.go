package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (cache lruCache) Set(key Key, value interface{}) bool {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	existedListItem, ok := cache.items[key]

	if ok {
		cache.queue.Remove(existedListItem)
	} else {
		length := cache.queue.Len()

		if length == cache.capacity {
			// remove the last item from the map and queue
			lastItem := cache.queue.Back()
			itemKey := lastItem.Value.(cacheItem).key
			delete(cache.items, itemKey)
			cache.queue.Remove(lastItem)
		}
	}

	newListItem := cacheItem{key, value}
	// add item to the front
	cache.queue.PushFront(newListItem)
	cache.items[key] = cache.queue.Front()

	return ok
}

func (cache lruCache) Get(key Key) (interface{}, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	var result interface{}

	value, ok := cache.items[key]

	if ok {
		cache.queue.MoveToFront(cache.items[key])
		result = value.Value.(cacheItem).value
	} else {
		result = nil
	}

	return result, ok
}

func (cache *lruCache) Clear() {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}