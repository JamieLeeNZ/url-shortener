package store

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
		return original, true
	}
	original, found := s.db.GetOriginalFromKey(key)
	if found {
		s.cache.Set(key, original)
	}
	return original, found
}

func (s *CachedStore) GetKeyFromOriginal(original string) (string, bool) {
	if key, found := s.cache.GetKeyFromOriginal(original); found {
		return key, true
	}
	key, found := s.db.GetKeyFromOriginal(original)
	if found {
		s.cache.Set(key, original)
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
