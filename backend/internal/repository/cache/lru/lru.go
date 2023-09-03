package lru

import (
	"container/list"
	"sync"
)

type LRUCache[T comparable] struct {
	capacity int
	l        *list.List
	cache    map[T]*list.Element
	lock     sync.Locker
}

type pair[T comparable] struct {
	Key   T
	Value any
}

func NewLRUCache[T comparable](cap int) *LRUCache[T] {
	return &LRUCache[T]{
		capacity: cap,
		l:        list.New(),
		cache:    make(map[T]*list.Element),
		lock:     &sync.Mutex{},
	}
}

func (l *LRUCache[T]) Put(k T, v any) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if elem, ok := l.cache[k]; ok {
		l.l.MoveToFront(elem)
		elem.Value = pair[T]{Key: k, Value: v}
	} else {
		if l.l.Len() >= l.capacity {
			delete(l.cache, l.l.Back().Value.(pair[T]).Key)
			l.l.Remove(l.l.Back())
		}
		l.cache[k] = l.l.PushFront(pair[T]{Key: k, Value: v})
	}
}

func (l *LRUCache[T]) Get(k T) (bool, any) {
	l.lock.Lock()
	defer l.lock.Unlock()

	elem, ok := l.cache[k]
	if !ok {
		var zero any
		return false, zero
	}

	l.l.MoveToFront(elem)
	return true, elem.Value.(pair[T]).Value
}
