package lru

import (
	"container/list"
	"errors"
	"fmt"
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
	size  uintptr
}

// Set stores this value in the interface with the given key.
// If storing this entry increases the size of the cache greater
// than the capacity, then the last used element is removed
func (l *LRU) Set(key interface{}, value interface{}) error {

	entrySize := reflect.TypeOf(value).Size()
	if entrySize > l.convenienceCap {
		return errors.New("Cannot store item that is larger than the cache capacity")
	}

	newSize := l.setSize(entrySize)

	_, ok := l.items[key]
	if !ok {
		entry := entry{
			key:   key,
			value: value,
			size:  entrySize,
		}
		listEntry := l.used.PushFront(entry)
		l.size = newSize
		l.items[key] = listEntry
	}

	return nil

}

func (l *LRU) setSize(entrySize uintptr) uintptr {

	newSize := l.size + entrySize
	if newSize <= l.convenienceCap {
		return newSize
	}

	l.removeLastElement()
	return l.setSize(entrySize)

}

func (l *LRU) removeLastElement() {

	lastElement := l.used.Back()
	l.used.Remove(lastElement)

	entry := lastElement.Value.(entry)
	delete(l.items, entry.key)
	l.size -= entry.size

}

// Get returns the element with this key, if it exists.
// The returned element becomes the most recently used
// element in the cache and is less likely to be removed
func (l *LRU) Get(key interface{}) (interface{}, bool) {

	entry, ok := l.get(key)
	if ok {
		l.used.MoveToFront(entry)
	}
	return entry, ok

}

// Peek returns the element with this key, if it exists.
// Unlike Get, the returned element does not become the
// most recently used element
func (l *LRU) Peek(key interface{}) (interface{}, bool) {
	return l.get(key)
}

func (l *LRU) get(key interface{}) (*list.Element, bool) {

	entry, ok := l.items[key]
	if !ok {
		return nil, ok
	}
	return entry, ok

}

// Del removes the entry with the given key. If they entry
// does not exist, an error is returned
func (l *LRU) Del(key interface{}) error {

	entry, ok := l.get(key)
	if !ok {
		return fmt.Errorf("cannot delete non-existent entry: %v", key)
	}
	l.removeEntry(entry)
	return nil

}

func (l *LRU) removeEntry(e *list.Element) {
	l.used.Remove(e)
	l.size -= e.Value.(entry).size
}

// Size returns the current size of the cache as an integer
func (l *LRU) Size() int {
	return int(l.size)
}

// Has returns true if the give key contains an entry in
// the cache
func (l *LRU) Has(key interface{}) bool {
	_, ok := l.items[key]
	return ok
}

// Reset clears the cache, throwing away all elements
func (l *LRU) Reset() {
	l.used.Init()
	l.items = map[interface{}]*list.Element{}
	l.size = uintptr(0)
}
