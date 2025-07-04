package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nerfthisdev/go-backend-test-task/internal/domain"
	"github.com/pkg/errors"
)

type Repository struct {
	DB *pgxpool.Pool
}

func Init(ctx context.Context) (*Repository, error) {
	dburi := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	cfg, err := pgxpool.ParseConfig(dburi)

	cfg.MaxConns = 25
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.MaxConnLifetime = 2 * time.Hour

	if err != nil {
		return nil, err
	}

	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Repository{
		DB: dbpool,
	}, nil
}

func (r *Repository) StoreRefreshToken(ctx context.Context, token domain.RefreshToken) error {
	query := `
	INSERT INTO refresh_tokens
	(guid, token_hash, session_id, user_agent, ip_address, created_at, expires_at)
	`
	_, err := r.DB.Exec(ctx, query, token.GUID, token.TokenHash, token.SessionID, token.UserAgent, token.IP, token.CreatedAt, token.ExpiresAt)
	if err != nil {
		return errors.Wrap(err, "failed to store token")
	}
	return nil
}

func (r *Repository) GetRefreshToken(ctx context.Context, guid string) (domain.RefreshToken, error) {
	query := `
			SELECT token_hash, session_id, user_agent, ip_address, created_at, expires_at
			FROM refresh_tokens
			WHERE guid = $1
		`
	row := r.DB.QueryRow(ctx, query, guid)

	var token domain.RefreshToken
	guiduuid, err := uuid.Parse(guid)

	token.GUID = guiduuid

	if err != nil {
		return domain.RefreshToken{}, err
	}

	err = row.Scan(&token.TokenHash, &token.SessionID, &token.UserAgent, &token.IP, &token.CreatedAt, &token.ExpiresAt)

	if err != nil {
		return domain.RefreshToken{}, err
	}

	return token, nil
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, guid string) error {
	query := `
			DELETE FROM refresh_tokens
			WHERE guid = $1
		`

	_, err := r.DB.Exec(ctx, query, guid)
	return err
}
