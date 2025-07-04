package auth

import (
	"context"
	"errors"
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

func NewAuthService(repo domain.TokenRepository, tokens domain.TokenService, users domain.UserRepository, logger *zap.Logger) *AuthService {
	return &AuthService{
		repo:   repo,
		tokens: tokens,
		users:  users,
		logger: logger,
	}
}

func (s *AuthService) Authorize(ctx context.Context, guid, useragent, ip string) (domain.TokenPair, error) {

	sessionID := uuid.NewString()

	if guid == "" {
		newGuid := uuid.New().String()
		err := s.users.CreateUser(ctx, newGuid)
		if err != nil {
			return domain.TokenPair{}, err
		}
		guid = newGuid
	}

	exists, err := s.users.UserExists(ctx, guid)

	if err != nil {
		s.logger.Error("failed to check user existance", zap.String("reason", err.Error()))
		return domain.TokenPair{}, err
	}

	if !exists {
		s.logger.Warn("unauthorized attempt with unknown guid", zap.String("guid", guid))
		return domain.TokenPair{}, errors.New("user does not exist")
	}

	accessToken, err := s.tokens.GenerateAccessToken(guid, sessionID)

	if err != nil {
		s.logger.Error("failed to generate access token", zap.String("reason", err.Error()))
		return domain.TokenPair{}, err
	}

	refreshPlain, err := s.tokens.GenerateRefreshToken()
	if err != nil {
		s.logger.Error("failed to generate refresh token", zap.String("reason", err.Error()))
		return domain.TokenPair{}, err
	}

	hashed, err := s.tokens.HashRefreshToken(refreshPlain)
	if err != nil {
		s.logger.Error("failed to hash refresh token", zap.String("reason", err.Error()))
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
