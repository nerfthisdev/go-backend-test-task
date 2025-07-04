package auth

import "time"

type JWTService struct {
	secret     string
	expiration time.Duration
}

func NewJwtService(secret string, expiration time.Duration) *JWTService {
	return &JWTService{secret: secret, expiration: expiration}
}
