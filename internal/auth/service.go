package auth

import (
	"errors"
	"time"

	"github.com/nerfthisdev/go-backend-test-task/internal/repository"
)

var ErrTokenNotFound = errors.New("token not found")

type AuthService struct {
	repo       *repository.Repository
	signingKey []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(repo *repository.Repository, signingKey string, accessTTL, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		repo:       repo,
		signingKey: []byte(signingKey),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}
