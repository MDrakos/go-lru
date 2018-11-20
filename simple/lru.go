package simple

import (
	"container/list"
	"errors"
	"reflect"
)

// LRU is a basic implmentation of a least
// recently used cache using a map
type LRU struct {
	capaticty      int
	convenienceCap uintptr
	size           uintptr
	used           *list.List
	items          map[interface{}]*list.Element
}

// NewLRU returns a new LRU capable of holding items
// up to the given capacity in bytes
// TODO: Consider setting a max
func NewLRU(capacity int) *LRU {
	return &LRU{
		capaticty:      capacity,
		convenienceCap: uintptr(capacity),
		size:           uintptr(0),
		used:           list.New(),
		items:          map[interface{}]*list.Element{},
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

	entrySize := reflect.TypeOf(value).Size()
	if entrySize > l.convenienceCap {
		return errors.New("Cannot store item that is larger than the cache capacity")
	}

	newSize := l.size + entrySize
	if newSize > l.convenienceCap {
		l.removeElement()
	}

	_, ok := l.items[key]
	if !ok {
		entry := entry{key: key, value: value}
		listEntry := l.used.PushFront(entry)
		l.size += newSize
		l.items[key] = listEntry
	}

	return nil

}

func (l *LRU) removeElement() {

	lastElement := l.used.Back()
	l.used.Remove(lastElement)

	entry := lastElement.Value.(entry)
	delete(l.items, entry.key)

	lastElementSize := reflect.TypeOf(lastElement).Size()
	l.size -= lastElementSize

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
