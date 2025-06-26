package store

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	db *pgxpool.Pool
}

func NewPostgresStore(connString string) (*PostgresStore, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}
	return &PostgresStore{db: pool}, nil
}

var _ URLStore = (*PostgresStore)(nil)

var mu sync.Mutex

func (s *PostgresStore) Set(key, originalURL string) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := s.db.Exec(context.Background(), `
		INSERT INTO url_mappings (key, original_url) VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET original_url = EXCLUDED.original_url
	`, key, originalURL)
	return err
}

func (s *PostgresStore) GetOriginalFromKey(key string) (string, bool) {
	var original string
	err := s.db.QueryRow(context.Background(),
		`SELECT original_url FROM url_mappings WHERE key = $1`, key).Scan(&original)
	if err != nil {
		return "", false
	}
	return original, true
}

func (s *PostgresStore) GetKeyFromOriginal(original string) (string, bool) {
	var key string
	err := s.db.QueryRow(context.Background(),
		`SELECT key FROM url_mappings WHERE original_url = $1`, original).Scan(&key)
	if err != nil {
		return "", false
	}
	return key, true
}

func (s *PostgresStore) ContainsKey(key string) bool {
	var exists bool
	err := s.db.QueryRow(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM url_mappings WHERE key = $1)`, key).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (s *PostgresStore) Update(key, newValue string) bool {
	cmdTag, err := s.db.Exec(context.Background(),
		`UPDATE url_mappings SET original_url = $1 WHERE key = $2`, newValue, key)
	if err != nil {
		return false
	}
	return cmdTag.RowsAffected() > 0
}

func (s *PostgresStore) Delete(key string) bool {
	cmdTag, err := s.db.Exec(context.Background(),
		`DELETE FROM url_mappings WHERE key = $1`, key)
	if err != nil {
		return false
	}
	return cmdTag.RowsAffected() > 0
}

func (s *PostgresStore) Close() error {
	if s.db != nil {
		s.db.Close()
	}
	return nil
}
