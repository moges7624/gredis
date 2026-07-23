package store

import (
	"sync"
	"time"
)

type entry struct {
	value     any
	expiresAt time.Time
}

func (e entry) expired(now time.Time) bool {
	return !e.expiresAt.IsZero() && now.After(e.expiresAt)
}

type Store struct {
	data map[string]entry
	mu   sync.RWMutex

	stopChan  chan struct{}
	cleanOnce sync.Once
}

func NewStore() *Store {
	s := &Store{
		data:     make(map[string]entry),
		stopChan: make(chan struct{}),
	}

	go s.cleanupLoop(10 * time.Second)
	return s
}

func (s *Store) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case now := <-ticker.C:
			s.cleanupExpired(now)
		}
	}
}

func (s *Store) cleanupExpired(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, e := range s.data {
		if e.expired(now) {
			delete(s.data, k)
		}
	}
}

func (s *Store) TTL(key string) (ttl time.Duration, exists bool, hasTTL bool) {
	e, exists := s.get(key)
	if !exists {
		return 0, false, false
	}

	if e.expiresAt.IsZero() {
		return 0, true, false
	}

	remaining := time.Until(e.expiresAt)
	if remaining < 0 {
		remaining = 0
	}

	return remaining, true, true
}

func (s *Store) get(key string) (entry, bool) {
	s.mu.RLock()
	e, exists := s.data[key]
	s.mu.RUnlock()

	if !exists {
		return entry{}, false
	}

	if !e.expired(time.Now()) {
		return e, true
	}

	s.mu.Lock()
	if e2, exists := s.data[key]; exists && e2.expired(time.Now()) {
		delete(s.data, key)
	}
	s.mu.Unlock()

	return entry{}, false
}

func (s *Store) Exists(keys []string) int {
	count := 0
	for _, k := range keys {
		if _, ok := s.get(k); ok {
			count++
		}
	}

	return count
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

func (s *Store) Close() {
	s.cleanOnce.Do(func() {
		close(s.stopChan)
	})
}
