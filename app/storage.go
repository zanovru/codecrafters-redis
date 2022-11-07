package main

import "sync"

type Storage struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (s *Storage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return "", false
	}
	return value, true
}

func (s *Storage) Set(key, value string) {
	s.mu.Lock()
	s.data[key] = value
	s.mu.Unlock()
}
