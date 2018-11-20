package lru

import (
	"container/list"
	"testing"

	"github.com/go-test/deep"
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

func TestLRU_Peek(t *testing.T) {

	l := NewLRU(42)

	err := l.Set("first", "first value")
	if err != nil {
		t.Fatal(err)
	}

	err = l.Set("second", "second value")
	if err != nil {
		t.Fatal(err)
	}

	secondEntry, ok := l.Peek("second")
	if !ok {
		t.Error("Expected to peek second element")
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
	if diff := deep.Equal(l.used.Front(), secondEntry); diff != nil {
		t.Errorf("expected cache's front entry to be second entry, \ngot: %s", diff)
	}

}

func TestLRU_Del(t *testing.T) {

	l := NewLRU(16)

	err := l.Set("first", "first value")
	if err != nil {
		t.Fatal(err)
	}

	err = l.Del("first")
	if err != nil {
		t.Error("expected no error when deleting existing entry")
	}
	if l.Size() != 0 {
		t.Errorf("expected cache size to be zero, got: %d", l.Size())
	}

}

func TestLRU_Del_NonExistent(t *testing.T) {

	l := NewLRU(16)

	err := l.Set("first", "first value")
	if err != nil {
		t.Fatal(err)
	}

	err = l.Del("second")
	if err == nil {
		t.Error("expected an error when attempting to delete a non-existent entry")
	}

}

func TestLRU_Has(t *testing.T) {

	l := NewLRU(8)

	err := l.Set("first", 1)
	if err != nil {
		t.Fatal(err)
	}

	ok := l.Has("first")
	if !ok {
		t.Error("expected Has to return true for existing entry")
	}

}

func TestLRU_Reset(t *testing.T) {

	l := NewLRU(16)

	err := l.Set("first", 1)
	if err != nil {
		t.Fatal(err)
	}

	err = l.Set("second", 2)
	if err != nil {
		t.Fatal(err)
	}

	l.Reset()
	if l.used.Len() != 0 {
		t.Errorf("exepcted Reset to clear list, got len: %d", l.used.Len())
	}
	if len(l.items) != 0 {
		t.Errorf("expected Reset to clear map, got len: %d", len(l.items))
	}
	if l.Size() != 0 {
		t.Errorf("expected Reset to set cache size to 0, got: %d", l.Size())
	}

}
