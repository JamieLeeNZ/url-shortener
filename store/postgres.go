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

func (s *PostgresStore) Set(ctx context.Context, key, originalURL string, userID string) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO url_mappings (key, original_url, user_id) VALUES ($1, $2, $3)
		ON CONFLICT (key) DO UPDATE SET original_url = EXCLUDED.original_url
	`, key, originalURL, userID)
	return err
}

func (s *PostgresStore) GetOriginalFromKey(ctx context.Context, key string) (string, string, bool) {
	var original, userID string
	err := s.db.QueryRow(ctx,
		`SELECT original_url, user_id FROM url_mappings WHERE key = $1`, key).Scan(&original, &userID)
	if err != nil {
		return "", "", false
	}
	return original, userID, true
}

func (s *PostgresStore) GetKeyFromOriginal(ctx context.Context, original string) (string, string, bool) {
	var key, userID string
	err := s.db.QueryRow(ctx,
		`SELECT key, user_id FROM url_mappings WHERE original_url = $1`, original).Scan(&key, &userID)
	if err != nil {
		return "", "", false
	}
	return key, userID, true
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

func (s *PostgresStore) GetURLsByUserID(ctx context.Context, userID string) ([]models.URLMapping, error) {
	rows, err := s.db.Query(ctx, `
		SELECT key, original_url, created_at
		FROM url_mappings
		WHERE user_id = $1
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []models.URLMapping
	for rows.Next() {
		var u models.URLMapping
		if err := rows.Scan(&u.Key, &u.Original, &u.CreatedAt); err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}

	return urls, nil
}
