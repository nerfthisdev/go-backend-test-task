package domain

import "context"

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, token RefreshToken) error
	GetRefreshToken(ctx context.Context, guid string) (RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, guid string) error
}

type UserRepository interface {
	UserExists(ctx context.Context, guid string) (bool, error)
	CreateUser(ctx context.Context, guid string) error
}
