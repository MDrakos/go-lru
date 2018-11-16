package simple

import (
	"container/list"
)

// LRU is a basic implmentation of a least
// recently used cache using a map
type LRU struct {
	capaticty int
	used      *list.List
	items     map[interface{}]*list.Element
}

// NewLRU returns a new LRU capable of holding items
// up to the given capacity in bytes
// TODO: Consider setting a max
func NewLRU(capacity int) *LRU {
	return &LRU{
		capaticty: capacity,
		used:      list.New(),
		items:     make(map[interface{}]*list.Element, capacity),
	}
}

type entry struct {
	key   interface{}
	value interface{}
}

// Set stores this value in the interface with the given key.
// If storing this entry increases the size of the cache greater
// than the capacity, then the last used element is removed
func (l *LRU) Set(key interface{}, value interface{}) error {

	if l.used.Len() >= l.capaticty {
		last := l.used.Back()
		l.used.Remove(last)
		entry := last.Value.(entry)
		delete(l.items, entry.key)
	}

	_, ok := l.items[key]
	if !ok {
		entry := entry{key: key, value: value}
		listEntry := l.used.PushFront(entry)
		l.items[key] = listEntry
	}
	return nil

}

// Get returns the element with this key, if it exists.
// The returned element becomes the most recently used
// element in the cache and is less likely to be removed
func (l *LRU) Get(key interface{}) (interface{}, bool) {

	entry, ok := l.items[key]
	if !ok {
		return nil, ok
	}
	l.used.MoveToFront(entry)
	return entry, ok

}
