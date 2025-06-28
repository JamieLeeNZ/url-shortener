package store

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
	ttl    time.Duration
}

type cachedURL struct {
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
}

var _ RedisClientProvider = (*RedisStore)(nil)
var _ URLStore = (*RedisStore)(nil)

func (r *RedisStore) RawClient() *redis.Client {
	return r.client
}

func NewRedisStore(addr, password string, db int, ttl time.Duration) (*RedisStore, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:      addr,
		Password:  password,
		DB:        db,
		TLSConfig: &tls.Config{},
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisStore{
		client: rdb,
		ttl:    ttl,
	}, nil
}

func (r *RedisStore) Set(ctx context.Context, key, original, userID string) error {
	data := cachedURL{OriginalURL: original, UserID: userID}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, key, jsonData, r.ttl).Err()
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, "original:"+original, key, r.ttl).Err()
	return err
}

func (r *RedisStore) GetOriginalFromKey(ctx context.Context, key string) (string, string, bool) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		return "", "", false
	}

	var data cachedURL
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return "", "", false
	}

	return data.OriginalURL, data.UserID, true
}

func (r *RedisStore) Update(ctx context.Context, key, newValue string) bool {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		return false
	}

	var data cachedURL
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return false
	}

	r.client.Del(ctx, "original:"+data.OriginalURL)

	data.OriginalURL = newValue

	newJSON, err := json.Marshal(data)
	if err != nil {
		return false
	}

	err = r.client.Set(ctx, key, newJSON, r.ttl).Err()
	if err != nil {
		return false
	}

	err = r.client.Set(ctx, "original:"+newValue, key, r.ttl).Err()
	return err == nil
}

func (r *RedisStore) ContainsKey(ctx context.Context, key string) bool {
	_, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		return false
	}
	return true
}

func (r *RedisStore) GetKeyFromOriginal(ctx context.Context, original string) (string, string, bool) {
	key, err := r.client.Get(ctx, "original:"+original).Result()
	if err == redis.Nil || err != nil {
		return "", "", false
	}

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		return "", "", false
	}

	var data cachedURL
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return "", "", false
	}

	return key, data.UserID, true
}

func (r *RedisStore) Delete(ctx context.Context, key string) bool {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		return false
	}

	var data cachedURL
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return false
	}

	err = r.client.Del(ctx, key).Err()
	if err != nil {
		return false
	}

	err = r.client.Del(ctx, "original:"+data.OriginalURL).Err()
	return err == nil
}

func (r *RedisStore) Close() error {
	return r.client.Close()
}
