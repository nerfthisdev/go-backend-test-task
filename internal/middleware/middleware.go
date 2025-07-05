package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/nerfthisdev/go-backend-test-task/internal/domain"
	"go.uber.org/zap"
)

type contextKey string

const (
	ContextUserGUIDKey    contextKey = "user_guid"
	ContextSessionIDKey   contextKey = "session_id"
	ContextAccessTokenKey contextKey = "access_token"
)

func Auth(logger *zap.Logger, tokens domain.TokenService, repo domain.TokenRepository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Warn("missing or invalid auth header")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := tokens.ValidateAccessToken(accessToken)
		if err != nil {
			logger.Warn("invalid token", zap.Error(err))
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		guidStr, okGUID := claims["sub"].(string)
		sessionID, okSessionID := claims["jti"].(string)

		parsedGUID, err := uuid.Parse(guidStr)
		if err != nil {
			logger.Warn("invalid guid in token", zap.Error(err))
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		logger.Warn("trying to auth", zap.String("guid", guidStr))
		logger.Warn("session id", zap.String("sessionID", sessionID))

		if !okGUID || !okSessionID {
			logger.Warn("missing token claims")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserGUIDKey, parsedGUID)
		ctx = context.WithValue(ctx, ContextSessionIDKey, sessionID)
		ctx = context.WithValue(ctx, ContextAccessTokenKey, accessToken)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
