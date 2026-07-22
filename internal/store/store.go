package store

import (
	"fmt"
	"strconv"
	"sync"
)

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

func (s *Store) Set(key string, val string) {
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

func (s *Store) IncrBy(key string, amount int64) (int64, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, exists := s.data[key]
	if !exists {
		return 0, false, nil
	}

	str, isStr := val.(string)
	if !isStr {
		return 0, true, fmt.Errorf("val is not a string")
	}

	intVal, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, true, err
	}

	newAmount := intVal + amount
	s.data[key] = strconv.FormatInt(newAmount, 10)

	return newAmount, true, nil
}
