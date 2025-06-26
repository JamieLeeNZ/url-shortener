package store

import (
	"log"
)

type CachedStore struct {
	cache URLStore
	db    URLStore
}

func NewCachedStore(cache, db URLStore) (*CachedStore, error) {
	return &CachedStore{cache: cache, db: db}, nil
}

var _ URLStore = (*CachedStore)(nil)

func (s *CachedStore) Set(key, originalURL string) error {
	if err := s.db.Set(key, originalURL); err != nil {
		return err
	}
	return s.cache.Set(key, originalURL)
}

func (s *CachedStore) GetOriginalFromKey(key string) (string, bool) {
	if original, found := s.cache.GetOriginalFromKey(key); found {
		log.Printf("[cache] hit for key: %s", key)
		return original, true
	}
	log.Printf("[cache] miss for key: %s", key)

	original, found := s.db.GetOriginalFromKey(key)
	if found {
		log.Printf("[db] fetched and caching key: %s", key)
		s.cache.Set(key, original)
	} else {
		log.Printf("[db] key not found: %s", key)
	}
	return original, found
}

func (s *CachedStore) GetKeyFromOriginal(original string) (string, bool) {
	if key, found := s.cache.GetKeyFromOriginal(original); found {
		log.Printf("[cache] hit for original URL: %s", original)
		return key, true
	}
	log.Printf("[cache] miss for original URL: %s", original)

	key, found := s.db.GetKeyFromOriginal(original)
	if found {
		log.Printf("[db] fetched and caching original URL: %s", original)
		s.cache.Set(key, original)
	} else {
		log.Printf("[db] original URL not found: %s", original)
	}
	return key, found
}

func (s *CachedStore) ContainsKey(key string) bool {
	if exists := s.cache.ContainsKey(key); exists {
		return true
	}
	return s.db.ContainsKey(key)
}

func (s *CachedStore) Update(key, newValue string) bool {
	ok := s.db.Update(key, newValue)
	if ok {
		s.cache.Update(key, newValue)
	}
	return ok
}

func (s *CachedStore) Delete(key string) bool {
	ok := s.db.Delete(key)
	if ok {
		s.cache.Delete(key)
	}
	return ok
}

func (c *CachedStore) Close() error {
	errDB := c.db.Close()
	errCache := c.cache.Close()

	if errDB != nil {
		return errDB
	}
	return errCache
}
