package domain

import "context"

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, token RefreshToken) error
	GetRefreshToken(ctx context.Context, guid string) (RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, guid string) error
}

type TokenService interface {
	GenerateAccessToken(guid, sessionID string) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateAccessToken(token string) (map[string]any, error)

	HashRefreshToken(token string) (string, error)
	CompareRefreshToken(token string, hash string) bool

	EncodeBase64(token string) string
	DecodeBase64(encoded string) (string, error)
}
