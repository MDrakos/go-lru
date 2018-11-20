package lru

import (
	"container/list"
	"testing"
)

func TestLRU_Set(t *testing.T) {

	l := NewLRU(100)
	err := l.Set("Test", "test string")
	if err != nil {
		t.Errorf("expected no error when calling set, got %s", err)
	}

}

func TestLRU_Get(t *testing.T) {

	// Each entry string is 16 bytes in size
	// We only want the cache to be large enough
	// for 2 items
	l := NewLRU(42)

	err := l.Set("first", "first value")
	if err != nil {
		t.Fatal(err)
	}

	err = l.Set("second", "second value")
	if err != nil {
		t.Fatal(err)
	}

	err = l.Set("third", "third value")
	if err != nil {
		t.Fatal(err)
	}

	_, ok := l.Get("first")
	if ok {
		t.Error("should not be able to get first element")
	}

	secondEntry, ok := l.Get("second")
	if !ok {
		t.Error("Expected to get second element back")
	}
	listElement, ok := secondEntry.(*list.Element)
	if !ok {
		t.Fatal("expected list entry type")
	}
	listEntry, ok := listElement.Value.(entry)
	if !ok {
		t.Fatal("expected list element to be an entry")
	}
	if listEntry.value != "second value" {
		t.Errorf("expected 'second value' got: %s", listEntry.value)
	}

}

func TestLRU_Set_TooLarge(t *testing.T) {

	l := NewLRU(10)

	err := l.Set("Big", "This is too big to store!")
	if err == nil {
		t.Error("expected error when setting an entry to large for the cache")
	}

}
