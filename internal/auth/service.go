package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nerfthisdev/go-backend-test-task/internal/domain"
	"go.uber.org/zap"
)

type AuthService struct {
	repo   domain.TokenRepository
	tokens domain.TokenService
	users  domain.UserRepository
	logger *zap.Logger
}

func NewAuthService(repo domain.TokenRepository, tokens domain.TokenService, logger *zap.Logger) *AuthService {
	return &AuthService{
		repo:   repo,
		tokens: tokens,
		logger: logger,
	}
}

func (s *AuthService) Authorize(ctx context.Context, guid, useragent, ip string) (domain.TokenPair, error) {

	sessionID := uuid.NewString()

	accessToken, err := s.tokens.GenerateAccessToken(guid, sessionID)

	if err != nil {
		s.logger.Error("failed to generate refresh token", zap.Error(err))
		return domain.TokenPair{}, err
	}

	refreshPlain, err := s.tokens.GenerateRefreshToken()
	if err != nil {
		s.logger.Error("failed to generate refresh token", zap.Error(err))
		return domain.TokenPair{}, err
	}

	hashed, err := s.tokens.HashRefreshToken(refreshPlain)
	if err != nil {
		s.logger.Error("failed to hash refresh token", zap.Error(err))
		return domain.TokenPair{}, err
	}

	guiduuid, err := uuid.Parse(guid)

	if err != nil {
		s.logger.Error("failed to parse uuid", zap.Error(err))
		return domain.TokenPair{}, err
	}

	refreshToken := domain.RefreshToken{
		GUID:      guiduuid,
		TokenHash: hashed,
		SessionID: sessionID,
		UserAgent: useragent,
		IP:        ip,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.repo.StoreRefreshToken(ctx, refreshToken); err != nil {
		s.logger.Error("failed to store refresh token", zap.Error(err))
		return domain.TokenPair{}, err
	}
	return domain.TokenPair{AccessToken: accessToken,
		RefreshToken: s.tokens.EncodeBase64(refreshPlain)}, nil

}
