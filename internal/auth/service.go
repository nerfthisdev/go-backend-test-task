package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nerfthisdev/go-backend-test-task/internal/model"
	"github.com/nerfthisdev/go-backend-test-task/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrTokenNotFound = errors.New("token not found")

type AuthService struct {
	Repo       *repository.Repository
	signingKey []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(repo *repository.Repository, signingKey string, accessTTL, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		Repo:       repo,
		signingKey: []byte(signingKey),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *AuthService) CreateTokenPair(ctx context.Context, guid uuid.UUID) (model.TokenPair, error) {
	// HMAC and SHA512
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		Subject:   guid.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	accessTokenString, err := accessToken.SignedString(s.signingKey)
	if err != nil {
		return model.TokenPair{}, err
	}

	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return model.TokenPair{}, err
	}

	refreshTokenString := base64.StdEncoding.EncodeToString(refreshTokenBytes)
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshTokenString), bcrypt.DefaultCost)
	if err != nil {
		return model.TokenPair{}, err
	}

	return model.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: string(hashedRefreshToken),
		IssuedAt:     time.Now().String(),
		RefreshExp:   time.Now().Add(s.refreshTTL).String(),
	}, err
}

func (s *AuthService) GetUserByGuid(ctx context.Context, guid uuid.UUID) (model.UserResponse, error) {
	userResponse, err := s.repo.SelectUser(ctx, guid)
	if err != nil {
		return model.UserResponse{}, err
	}

	return *userResponse, nil
}
