package lru

// LRU defines the provided API methods for the
// least recently used cache implementations
type LRU interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Peek(key string) (interface{}, error)
	Delete(key string) error
	Reset() error
	Has() (bool, error)
}
