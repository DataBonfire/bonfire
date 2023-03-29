package cache

import "sync"

type node struct {
	key   interface{}
	value interface{}
	prev  *node
	next  *node
}

type LRUCache struct {
	capacity int
	m        map[interface{}]*node
	head     *node
	tail     *node
}

func (cache *LRUCache) Delete(key interface{}) {
	n, ok := cache.m[key]
	if ok {
		cache.remove(n, false)
	}
}

func (cache *LRUCache) Len() int {
	return len(cache.m)
}

func NewLRUCache(capacity int) *LRUCache {
	cache := &LRUCache{}
	cache.capacity = capacity
	cache.m = map[interface{}]*node{}

	head := &node{}
	tail := &node{}

	head.next = tail
	tail.prev = head

	cache.head = head
	cache.tail = tail

	return cache
}

func (cache *LRUCache) remove(n *node, listOnly bool) {
	if !listOnly {
		delete(cache.m, n.key)
	}
	n.prev.next = n.next
	n.next.prev = n.prev
}

func (cache *LRUCache) add(n *node, listOnly bool) {
	if !listOnly {
		cache.m[n.key] = n
	}
	headNext := cache.head.next
	cache.head.next = n
	headNext.prev = n
	n.next = headNext
	n.prev = cache.head
}

func (cache *LRUCache) moveToHead(n *node) {
	cache.remove(n, true)
	cache.add(n, true)
}

func (cache *LRUCache) Get(key interface{}) (value interface{}, ok bool) {
	n, ok := cache.m[key]
	if ok {
		cache.moveToHead(n)
		return n.value, ok
	} else {
		return nil, ok
	}
}

func (cache *LRUCache) Put(key interface{}, value interface{}) {
	n, ok := cache.m[key]
	if ok {
		cache.remove(n, false)
	} else {
		n = &node{key, value, nil, nil}
		if len(cache.m) >= cache.capacity {
			cache.remove(cache.tail.prev, false)
		}
	}
	cache.add(n, false)
}

type SyncLRUCache struct {
	rwm      sync.RWMutex
	lruCache *LRUCache
}

func NewSyncLRUCache(capacity int) ICache {
	cache := &SyncLRUCache{}
	cache.lruCache = NewLRUCache(capacity)
	return cache
}

func (cache *SyncLRUCache) Get(key interface{}) (value interface{}, ok bool) {
	cache.rwm.Lock()
	defer cache.rwm.Unlock()
	return cache.lruCache.Get(key)
}

func (cache *SyncLRUCache) Put(key interface{}, value interface{}) {
	cache.rwm.Lock()
	defer cache.rwm.Unlock()
	cache.lruCache.Put(key, value)
}

func (cache *SyncLRUCache) Delete(key interface{}) {
	cache.rwm.Lock()
	defer cache.rwm.Unlock()
	cache.lruCache.Delete(key)
}

func (cache *SyncLRUCache) Len() int {
	cache.rwm.Lock()
	defer cache.rwm.Unlock()
	return cache.lruCache.Len()
}
