package auth

import (
	"github.com/nerfthisdev/go-backend-test-task/internal/domain"
	"go.uber.org/zap"
)

type AuthService struct {
	repo   domain.TokenRepository
	tokens domain.TokenService
	logger *zap.Logger
}

func NewAuthService(repo domain.TokenRepository, tokens domain.TokenService, logger *zap.Logger) *AuthService {
	return &AuthService{
		repo:   repo,
		tokens: tokens,
		logger: logger,
	}
}
