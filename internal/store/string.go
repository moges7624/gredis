package store

import (
	"errors"
	"strconv"
	"time"
)

var ErrNotInteger = errors.New(
	"value is not an integer or out of range",
)

type StringValue struct {
	data string
}

func (s *StringValue) Type() ValueType {
	return StringType
}

type SetOptions struct {
	TTL             time.Duration
	OnlyIfExists    bool
	OnlyIfNotExists bool
}

func (s *Store) Get(key string) (val string, exists bool, err error) {
	s.mu.RLock()
	e, exists := s.get(key)
	s.mu.RUnlock()

	if !exists {
		return "", false, nil
	}

	strVal, ok := e.value.(*StringValue)
	if !ok {
		return "", true, ErrWrongType
	}

	return strVal.data, true, nil
}

func (s *Store) Set(key, val string, opts SetOptions) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, exists := s.data[key]
	if exists && e.expired(time.Now()) {
		exists = false
	}

	if opts.OnlyIfExists && !exists {
		return false
	}

	if opts.OnlyIfNotExists && exists {
		return false
	}

	ne := entry{
		value: &StringValue{data: val},
	}
	if opts.TTL > 0 {
		ne.expiresAt = time.Now().Add(opts.TTL)
	}

	s.data[key] = ne
	return true
}

func (s *Store) IncrBy(key string, amount int64) (int64, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, exists := s.data[key]
	if !exists {
		return 0, false, nil
	}

	strVal, ok := e.value.(*StringValue)
	if !ok {
		return 0, true, ErrWrongType
	}

	intVal, err := strconv.ParseInt(strVal.data, 10, 64)
	if err != nil {
		return 0, true, ErrNotInteger
	}

	newAmount := intVal + amount
	s.data[key] = entry{
		value: &StringValue{strconv.FormatInt(newAmount, 10)},
	}

	return newAmount, true, nil
}
