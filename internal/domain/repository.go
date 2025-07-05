package domain

import (
	"context"

	"github.com/google/uuid"
)

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, token RefreshToken) error
	GetRefreshToken(ctx context.Context, guid uuid.UUID) (RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, guid uuid.UUID) error
	SessionExists(ctx context.Context, sessionID string) (bool, error)
}

type UserRepository interface {
	UserExists(ctx context.Context, guid uuid.UUID) (bool, error)
	CreateUser(ctx context.Context, guid uuid.UUID) error
}
