package store

import "sync"

type MemoryStore struct {
	mu            sync.RWMutex
	keyToOriginal map[string]string
	originalToKey map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		keyToOriginal: make(map[string]string),
		originalToKey: make(map[string]string),
	}
}

func (s *MemoryStore) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.keyToOriginal[key] = value
	s.originalToKey[value] = key
}

func (s *MemoryStore) GetOriginalFromKey(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.keyToOriginal[key]
	return value, ok
}

func (s *MemoryStore) GetKeyFromOriginal(original string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key, ok := s.originalToKey[original]
	return key, ok
}

func (s *MemoryStore) ContainsKey(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.keyToOriginal[key]
	return ok
}
