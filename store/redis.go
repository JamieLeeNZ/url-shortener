package store

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
	ttl    time.Duration
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

func (r *RedisStore) GetOriginalFromKey(key string) (string, bool) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false
	}
	if err != nil {
		return "", false
	}
	return val, true
}

func (r *RedisStore) GetKeyFromOriginal(original string) (string, bool) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, "original:"+original).Result()
	if err == redis.Nil {
		return "", false
	}
	if err != nil {
		return "", false
	}
	return val, true
}

func (r *RedisStore) Set(key, original string) error {
	ctx := context.Background()

	err := r.client.Set(ctx, key, original, r.ttl).Err()
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, "original:"+original, key, r.ttl).Err()
	return err
}

func (r *RedisStore) Update(key, newValue string) bool {
	ctx := context.Background()
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil || exists == 0 {
		return false
	}

	oldOriginal, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return false
	}

	r.client.Del(ctx, "original:"+oldOriginal)

	err = r.client.Set(ctx, key, newValue, r.ttl).Err()
	if err != nil {
		return false
	}

	err = r.client.Set(ctx, "original:"+newValue, key, r.ttl).Err()
	return err == nil
}

func (r *RedisStore) Delete(key string) bool {
	ctx := context.Background()
	original, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		return false
	}

	err = r.client.Del(ctx, key).Err()
	if err != nil {
		return false
	}

	err = r.client.Del(ctx, "original:"+original).Err()
	return err == nil
}

func (r *RedisStore) Close() error {
	return r.client.Close()
}
