package domain

import (
	"time"

	"github.com/google/uuid"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	GUID      uuid.UUID
	TokenHash string
	SessionID string
	IP        string
	UserAgent string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type TokenService interface {
	GenerateAccessToken(guid uuid.UUID, sessionID string) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateAccessToken(token string) (map[string]any, error)

	HashRefreshToken(token string) (string, error)
	CompareRefreshToken(token string, hash string) bool

	EncodeBase64(token string) string
	DecodeBase64(encoded string) (string, error)
}
