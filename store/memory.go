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

var _ URLStore = (*MemoryStore)(nil)

func (s *MemoryStore) Set(key, originalURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.keyToOriginal[key] = originalURL
	s.originalToKey[originalURL] = key
	return nil
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

func (s *MemoryStore) Update(key, newValue string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	original, ok := s.keyToOriginal[key]
	if !ok {
		return false
	}

	if existingKey, exists := s.originalToKey[newValue]; exists && existingKey != key {
		return false
	}

	if original == newValue {
		return true
	}

	delete(s.originalToKey, original)
	s.keyToOriginal[key] = newValue
	s.originalToKey[newValue] = key

	return true
}

func (s *MemoryStore) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if original, ok := s.keyToOriginal[key]; ok {
		delete(s.keyToOriginal, key)
		delete(s.originalToKey, original)
		return true
	}
	return false
}

func (s *MemoryStore) Close() error {
	// No resources to close in memory store
	return nil
}
