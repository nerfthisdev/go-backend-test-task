package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JWTService struct {
	secret    string
	accessTTL time.Duration
}

func NewJwtService(secret string, expiration time.Duration) *JWTService {
	return &JWTService{secret: secret, accessTTL: expiration}
}

func (s *JWTService) GenerateAccessToken(guid, sessionID string) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		Subject:   guid,
		ID:        sessionID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTTL)),
	})

	accesTokenString, err := accessToken.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}
	return accesTokenString, nil
}

func (s *JWTService) GenerateRefreshToken() (string, error) {
	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return "", err
	}
	return string(refreshTokenBytes), nil
}

func (s *JWTService) ValidateAccessToken(token string) (map[string]any, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !parsed.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}

func (s *JWTService) HashRefreshToken(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *JWTService) CompareRefreshToken(token string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	return err == nil
}

func (s *JWTService) EncodeBase64(token string) string {
	return base64.StdEncoding.EncodeToString([]byte(token))
}

func (s *JWTService) DecodeBase64(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
