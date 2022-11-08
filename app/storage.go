package main

import (
	"sync"
	"time"
)

type ValueWithExpiration struct {
	value string
	tms   time.Time
}

type Storage struct {
	mu   sync.RWMutex
	data map[string]ValueWithExpiration
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]ValueWithExpiration),
	}
}

func (s *Storage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return "", false
	}
	if isExpired(v) {
		delete(s.data, key)
		return "", false
	}
	return v.value, true
}

func (s *Storage) Set(key, value string) {
	s.mu.Lock()
	s.data[key] = ValueWithExpiration{
		value: value,
	}
	s.mu.Unlock()
}

func (s *Storage) SetWithExpiration(key, value string, duration time.Duration) {
	s.mu.Lock()
	s.data[key] = ValueWithExpiration{
		value: value,
		tms:   time.Now().Add(duration),
	}
	s.mu.Unlock()
}

func isExpired(v ValueWithExpiration) bool {
	if v.tms.IsZero() {
		return false
	}
	if v.tms.Before(time.Now()) {
		return true
	}
	return false
}
