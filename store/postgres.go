package store

import (
	"context"
	"time"

	"github.com/JamieLeeNZ/url-shortener/models"
	"github.com/jackc/pgx/v5"
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

	config.ConnConfig.StatementCacheCapacity = 0
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

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
var _ UserStore = (*PostgresStore)(nil)

func (s *PostgresStore) Set(ctx context.Context, key, originalURL string) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO url_mappings (key, original_url) VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET original_url = EXCLUDED.original_url
	`, key, originalURL)
	return err
}

func (s *PostgresStore) GetOriginalFromKey(ctx context.Context, key string) (string, bool) {
	var original string
	err := s.db.QueryRow(ctx,
		`SELECT original_url FROM url_mappings WHERE key = $1`, key).Scan(&original)
	if err != nil {
		return "", false
	}
	return original, true
}

func (s *PostgresStore) GetKeyFromOriginal(ctx context.Context, original string) (string, bool) {
	var key string
	err := s.db.QueryRow(ctx,
		`SELECT key FROM url_mappings WHERE original_url = $1`, original).Scan(&key)
	if err != nil {
		return "", false
	}
	return key, true
}

func (s *PostgresStore) ContainsKey(ctx context.Context, key string) bool {
	var exists bool
	err := s.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM url_mappings WHERE key = $1)`, key).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (s *PostgresStore) Update(ctx context.Context, key, newValue string) bool {
	cmdTag, err := s.db.Exec(ctx,
		`UPDATE url_mappings SET original_url = $1 WHERE key = $2`, newValue, key)
	if err != nil {
		return false
	}
	return cmdTag.RowsAffected() > 0
}

func (s *PostgresStore) Delete(ctx context.Context, key string) bool {
	cmdTag, err := s.db.Exec(ctx,
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

func (s *PostgresStore) GetOrCreateUser(ctx context.Context, user models.User) (models.User, error) {
	var existing models.User

	err := s.db.QueryRow(ctx, `
		SELECT id, email, name, picture_url, created_at
		FROM users WHERE id = $1`,
		user.ID,
	).Scan(&existing.ID, &existing.Email, &existing.Name, &existing.Picture, &existing.CreatedAt)

	if err == nil {
		return existing, nil
	}

	if err != pgx.ErrNoRows {
		return models.User{}, err
	}

	_, err = s.db.Exec(ctx, `
		INSERT INTO users (id, email, name, picture_url, created_at)
		VALUES ($1, $2, $3, $4, NOW())`,
		user.ID, user.Email, user.Name, user.Picture,
	)
	if err != nil {
		return models.User{}, err
	}

	user.CreatedAt = time.Now()
	return user, nil
}
