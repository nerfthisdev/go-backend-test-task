package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (s *AuthService) CreateTokenPair(ctx context.Context, guid uuid.UUID) (TokenPair, error) {
	// HMAC and SHA512
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		Subject:   guid.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	accessTokenString, err := accessToken.SignedString(s.signingKey)
	if err != nil {
		return TokenPair{}, err
	}

	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return TokenPair{}, err
	}

	refreshTokenString := base64.StdEncoding.EncodeToString(refreshTokenBytes)

	return TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, err
}
