package hw04lrucache

import "sync"

type Key string

type Cache[T any] interface {
	Set(key Key, value T) bool
	Get(key Key) (T, bool)
	Clear()
}

type cacheValue[T any] struct {
	key   Key
	value T
}

type lruCache[T any] struct {
	capacity int
	queue    List[cacheValue[T]]
	items    map[Key]*ListItem[cacheValue[T]]
	m        sync.Mutex
}

func (lc *lruCache[T]) Set(key Key, value T) bool {
	lc.m.Lock()
	defer lc.m.Unlock()
	i, ok := lc.items[key]
	if ok {
		i.Value = cacheValue[T]{key: key, value: value}
		lc.queue.MoveToFront(i)
	} else {
		switch {
		case lc.queue.Len() < lc.capacity:
			newItem := lc.queue.PushFront(cacheValue[T]{key: key, value: value})
			lc.items[key] = newItem
		case lc.queue.Len() == lc.capacity:
			delKey := lc.queue.Back().Value.key
			lc.queue.Remove(lc.queue.Back())
			delete(lc.items, delKey)
			newItem := lc.queue.PushFront(cacheValue[T]{key: key, value: value})
			lc.items[key] = newItem
		default:
			panic("queue length is bigger than cache capacity")
		}
	}
	return ok
}

func (lc *lruCache[T]) Get(key Key) (T, bool) {
	lc.m.Lock()
	defer lc.m.Unlock()
	var v T
	i, ok := lc.items[key]
	if ok {
		lc.queue.MoveToFront(i)
		v = i.Value.value
	}
	return v, ok
}

func (lc *lruCache[T]) Clear() {
	lc.m.Lock()
	defer lc.m.Unlock()
	lc.queue = NewList[cacheValue[T]]()
	clear(lc.items)
}

func NewCache[T any](capacity int) Cache[T] {
	return &lruCache[T]{
		capacity: capacity,
		queue:    NewList[cacheValue[T]](),
		items:    make(map[Key]*ListItem[cacheValue[T]], capacity),
	}
}
