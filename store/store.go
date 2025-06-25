package store

type URLStore interface {
	Set(key, originalURL string) error
	GetOriginalFromKey(key string) (string, bool)
	GetKeyFromOriginal(originalURL string) (string, bool)
	Update(key, newURL string) bool
	Delete(key string) bool
	ContainsKey(key string) bool
	Close() error
}
