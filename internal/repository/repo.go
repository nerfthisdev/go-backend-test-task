package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nerfthisdev/go-backend-test-task/internal/config"
	"github.com/nerfthisdev/go-backend-test-task/internal/domain"
	"github.com/pkg/errors"
)

type Repository struct {
	DB *pgxpool.Pool
}

func Init(ctx context.Context, cfg config.Config) (*Repository, error) {
	dburi := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBAddress,
		cfg.DBPort,
		cfg.DBName,
	)

	pgCfg, err := pgxpool.ParseConfig(dburi)
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}

	pgCfg.MaxConns = 25
	pgCfg.MaxConnIdleTime = 5 * time.Minute
	pgCfg.MaxConnLifetime = 2 * time.Hour

	dbpool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	return &Repository{DB: dbpool}, nil
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
