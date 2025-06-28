package store

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type CachedStore struct {
	cache URLStore
	db    URLStore
}

type RedisClientProvider interface {
	RawClient() *redis.Client
}

func NewCachedStore(cache, db URLStore) (*CachedStore, error) {
	return &CachedStore{cache: cache, db: db}, nil
}

var _ URLStore = (*CachedStore)(nil)

func (c *CachedStore) RedisClient() *redis.Client {
	if provider, ok := c.cache.(RedisClientProvider); ok {
		return provider.RawClient()
	}
	return nil
}

func (s *CachedStore) Set(ctx context.Context, key, originalURL string, userID string) error {
	if err := s.db.Set(ctx, key, originalURL, userID); err != nil {
		return err
	}
	return s.cache.Set(ctx, key, originalURL, userID)
}

func (s *CachedStore) GetOriginalFromKey(ctx context.Context, key string) (string, string, bool) {
	if original, userID, found := s.cache.GetOriginalFromKey(ctx, key); found {
		log.Printf("[cache] hit for key: %s", key)
		return original, userID, true
	}
	log.Printf("[cache] miss for key: %s", key)

	original, userID, found := s.db.GetOriginalFromKey(ctx, key)
	if found {
		log.Printf("[db] fetched and caching key: %s", key)
		s.cache.Set(ctx, key, original, userID)
	} else {
		log.Printf("[db] key not found: %s", key)
	}
	return original, userID, found
}

func (s *CachedStore) GetKeyFromOriginal(ctx context.Context, original string) (string, string, bool) {
	if key, userID, found := s.cache.GetKeyFromOriginal(ctx, original); found {
		log.Printf("[cache] hit for original URL: %s", original)
		return key, userID, true
	}
	log.Printf("[cache] miss for original URL: %s", original)

	key, userID, found := s.db.GetKeyFromOriginal(ctx, original)
	if found {
		log.Printf("[db] fetched and caching original URL: %s", original)
		s.cache.Set(ctx, key, original, userID)
	} else {
		log.Printf("[db] original URL not found: %s", original)
	}
	return key, userID, found
}

func (s *CachedStore) ContainsKey(ctx context.Context, key string) bool {
	if exists := s.cache.ContainsKey(ctx, key); exists {
		return true
	}
	return s.db.ContainsKey(ctx, key)
}

func (s *CachedStore) Update(ctx context.Context, key, newValue string) bool {
	ok := s.db.Update(ctx, key, newValue)
	if ok {
		s.cache.Update(ctx, key, newValue)
	}
	return ok
}

func (s *CachedStore) Delete(ctx context.Context, key string) bool {
	ok := s.db.Delete(ctx, key)
	if ok {
		s.cache.Delete(ctx, key)
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
