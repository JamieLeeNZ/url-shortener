package store

import (
	"context"
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

var _ URLStore = (*MemoryStore)(nil)

func (s *PostgresStore) Set(key, originalURL string) error {
	// TODO:
	return nil
}

func (s *PostgresStore) GetOriginalFromKey(key string) (string, bool) {
	// TODO:
	return "", false
}

func (s *PostgresStore) GetKeyFromOriginal(original string) (string, bool) {
	//  TODO:
	return "", false
}

func (s *PostgresStore) ContainsKey(key string) bool {
	// TODO:
	return false
}

func (s *PostgresStore) Update(key, newValue string) bool {
	// TODO:
	return false
}

func (s *PostgresStore) Delete(key string) bool {
	// TODO:
	return false
}

func (s *PostgresStore) Close() error {
	// TODO:
	return nil
}
