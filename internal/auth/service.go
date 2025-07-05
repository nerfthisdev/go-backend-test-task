package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

func (s *AuthService) Authorize(ctx context.Context, guid *uuid.UUID, useragent, ip string) (domain.TokenPair, error) {
	sessionID := uuid.NewString()

	if guid == nil {
		newGuid := uuid.New()
		err := s.users.CreateUser(ctx, newGuid)
		if err != nil {
			return domain.TokenPair{}, err
		}
		guid = &newGuid
	}

	exists, err := s.users.UserExists(ctx, *guid)
	if err != nil {
		s.logger.Error("failed to check user existance", zap.String("reason", err.Error()))
		return domain.TokenPair{}, err
	}

	if !exists {
		s.logger.Warn("unauthorized attempt with unknown guid", zap.String("guid", guid.String()))
		return domain.TokenPair{}, fmt.Errorf("user does not exist")
	}

	accessToken, err := s.tokens.GenerateAccessToken(*guid, sessionID)
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
		s.logger.Error("failed to hash refresh token", zap.Error(err))
		return domain.TokenPair{}, err
	}

	refreshToken := domain.RefreshToken{
		GUID:      *guid,
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
	return domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: s.tokens.EncodeBase64(refreshPlain),
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, guid uuid.UUID, sessionID, refreshToken, userAgent, ip string) (domain.TokenPair, error) {
	stored, err := s.repo.GetRefreshToken(ctx, guid)
	if err != nil {
		s.logger.Error("refresh failed: no stored token", zap.Error(err))
		return domain.TokenPair{}, nil
	}

	if stored.SessionID != sessionID {
		s.logger.Warn("session id doesnt match", zap.String("guid", guid.String()))
		s.repo.DeleteRefreshToken(ctx, guid)
		return domain.TokenPair{}, fmt.Errorf("unauthorized")
	}

	decoded, err := s.tokens.DecodeBase64(refreshToken)
	if err != nil || !s.tokens.CompareRefreshToken(decoded, stored.TokenHash) {
		s.logger.Warn("refresh token mismatch or tampered", zap.String("guid", guid.String()))
		s.repo.DeleteRefreshToken(ctx, guid)
		return domain.TokenPair{}, fmt.Errorf("unauthorized")
	}

	if stored.UserAgent != userAgent {
		s.logger.Warn("user-agent mismatch", zap.String("guid", guid.String()))
		s.repo.DeleteRefreshToken(ctx, guid)
		return domain.TokenPair{}, fmt.Errorf("unauthorized")
	}

	if stored.IP != ip {
		go s.sendIPChangeWebhook(guid, stored.IP, ip, userAgent)
	}

	return s.Authorize(ctx, &guid, userAgent, ip)
}

func (s *AuthService) sendIPChangeWebhook(guid uuid.UUID, oldIP, newIP, ua string) {
	type WebhookPayload struct {
		GUID      string `json:"guid"`
		OldIP     string `json:"old_ip"`
		NewIP     string `json:"new_ip"`
		UserAgent string `json:"user_agent"`
		Time      string `json:"time"`
	}

	payload := WebhookPayload{
		GUID:      guid.String(),
		OldIP:     oldIP,
		NewIP:     newIP,
		UserAgent: ua,
		Time:      time.Now().Format(time.RFC3339),
	}

	body, _ := json.Marshal(payload)
	_, err := http.Post(os.Getenv("WEBHOOK_URL"), "application/json", bytes.NewReader(body))
	if err != nil {
		s.logger.Error("failed to send IP change webhook", zap.Error(err))
	}
}
