package store

import "sync"

type Store struct {
	data map[string]any
	mu   sync.RWMutex
}

func NewStore() *Store {
	store := &Store{data: make(map[string]any)}
	return store
}

func (s *Store) Get(key string) (any, bool) {
	s.mu.RLock()
	val, exists := s.data[key]
	s.mu.RUnlock()

	return val, exists
}

func (s *Store) Set(key string, val any) {
	s.mu.Lock()
	s.data[key] = val
	s.mu.Unlock()
}
