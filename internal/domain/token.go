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
