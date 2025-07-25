package store

import "context"

type URLStore interface {
	Set(ctx context.Context, key, originalURL string, userID string) error
	GetOriginalFromKey(ctx context.Context, key string) (string, string, bool)
	GetKeyFromOriginal(ctx context.Context, original string) (string, string, bool)
	ContainsKey(ctx context.Context, key string) bool
	Update(ctx context.Context, key, newValue string) bool
	Delete(ctx context.Context, key string) bool
	Close() error
}
