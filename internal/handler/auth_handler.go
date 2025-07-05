package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/nerfthisdev/go-backend-test-task/internal/auth"
)

type RefreshRequest struct {
	GUID         uuid.UUID `json:"guid"`
	RefreshToken string    `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthHandler struct {
	auth *auth.AuthService
}

func NewAuthHandler(service *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		auth: service,
	}
}

func (h *AuthHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")

	userAgent := r.UserAgent()
	ip := r.RemoteAddr

	tokens, err := h.auth.Authorize(r.Context(), guid, userAgent, ip)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
