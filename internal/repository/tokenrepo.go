package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nerfthisdev/go-backend-test-task/internal/domain"
)

type TokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) StoreRefreshToken(ctx context.Context, token domain.RefreshToken) error {
	query := `

		INSERT INTO refresh_tokens
		(guid, token_hash, session_id, user_agent, ip_address, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (guid) DO UPDATE
		SET token_hash = EXCLUDED.token_hash,
			session_id = EXCLUDED.session_id,
			user_agent = EXCLUDED.user_agent,
			ip_address = EXCLUDED.ip_address,
			created_at = EXCLUDED.created_at,
			expires_at = EXCLUDED.expires_at
`
	_, err := r.db.Exec(ctx, query, token.GUID, token.TokenHash, token.SessionID, token.UserAgent, token.IP, token.CreatedAt, token.ExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}
	return nil
}

func (r *TokenRepository) GetRefreshToken(ctx context.Context, guid uuid.UUID) (domain.RefreshToken, error) {
	query := `
			SELECT token_hash, session_id, user_agent, ip_address, created_at, expires_at
			FROM refresh_tokens
			WHERE guid = $1
		`
	row := r.db.QueryRow(ctx, query, guid)

	var token domain.RefreshToken
	token.GUID = guid

	if err := row.Scan(&token.TokenHash, &token.SessionID, &token.UserAgent, &token.IP, &token.CreatedAt, &token.ExpiresAt); err != nil {
		return domain.RefreshToken{}, err
	}

	return token, nil
}

func (r *TokenRepository) DeleteRefreshToken(ctx context.Context, guid uuid.UUID) error {
	query := `
			DELETE FROM refresh_tokens
			WHERE guid = $1
		`

	_, err := r.db.Exec(ctx, query, guid)
	return err
}

func (r *TokenRepository) SessionExists(ctx context.Context, sessionID string) (bool, error) {
	query := `SELECT 1 FROM refresh_tokens WHERE session_id = $1 LIMIT 1`
	row := r.db.QueryRow(ctx, query, sessionID)

	var dummy int
	err := row.Scan(&dummy)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
