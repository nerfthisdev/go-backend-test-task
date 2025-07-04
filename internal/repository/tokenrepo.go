package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nerfthisdev/go-backend-test-task/internal/domain"
	"github.com/pkg/errors"
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
		return errors.Wrap(err, "failed to store token")
	}
	return nil
}

func (r *TokenRepository) GetRefreshToken(ctx context.Context, guid string) (domain.RefreshToken, error) {
	query := `
			SELECT token_hash, session_id, user_agent, ip_address, created_at, expires_at
			FROM refresh_tokens
			WHERE guid = $1
		`
	row := r.db.QueryRow(ctx, query, guid)

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

func (r *TokenRepository) DeleteRefreshToken(ctx context.Context, guid string) error {
	query := `
			DELETE FROM refresh_tokens
			WHERE guid = $1
		`

	_, err := r.db.Exec(ctx, query, guid)
	return err
}
