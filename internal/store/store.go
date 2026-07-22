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

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	val, exists := s.data[key]
	s.mu.RUnlock()

	str, isStr := val.(string)
	if !isStr {
		return "", false
	}

	return str, exists
}

func (s *Store) Set(key string, val any) {
	s.mu.Lock()
	s.data[key] = val
	s.mu.Unlock()
}

func (s *Store) Delete(key []string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := 0
	for _, k := range key {
		if _, exists := s.data[k]; exists {
			count++
		}
		delete(s.data, k)
	}

	return count
}
