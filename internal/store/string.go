package store

import (
	"fmt"
	"strconv"
	"time"
)

type SetOptions struct {
	TTL             time.Duration
	OnlyIfExists    bool
	OnlyIfNotExists bool
}

func (s *Store) Get(key string) (string, bool) {
	e, exists := s.get(key)
	if !exists {
		return "", false
	}

	str, isStr := e.value.(string)
	if !isStr {
		return "", false
	}

	return str, true
}

func (s *Store) Set(key, val string, opts SetOptions) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e := entry{value: val}
	if opts.TTL > 0 {
		e.expiresAt = time.Now().Add(opts.TTL)
	}

	s.data[key] = e
}

func (s *Store) IncrBy(key string, amount int64) (int64, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, exists := s.data[key]
	if !exists {
		return 0, false, nil
	}

	str, isStr := e.value.(string)
	if !isStr {
		return 0, true, fmt.Errorf("val is not a string")
	}

	intVal, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, true, err
	}

	newAmount := intVal + amount
	s.data[key] = entry{value: strconv.FormatInt(newAmount, 10)}

	return newAmount, true, nil
}
