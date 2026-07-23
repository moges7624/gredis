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
	e, exists := s.get(key)
	if !exists {
		return "", false, nil
	}

	strVal, ok := e.value.(*StringValue)
	if !ok {
		return "", true, ErrWrongType
	}

	return strVal.data, true, nil
}

func (s *Store) Set(key, val string, opts SetOptions) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e := entry{
		value: &StringValue{data: val},
	}
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
