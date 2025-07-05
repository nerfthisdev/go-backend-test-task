package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/nerfthisdev/go-backend-test-task/internal/auth"
	"github.com/nerfthisdev/go-backend-test-task/internal/middleware"
)

type RefreshRequest struct {
	GUID         uuid.UUID `json:"guid"`
	RefreshToken string    `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MeResponse struct {
	GUID string `json:"guid"`
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

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	sessionID := r.Context().Value(middleware.ContextSessionIDKey).(string)

	userAgent := r.UserAgent()
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.RemoteAddr
	}

	tokens, err := h.auth.Refresh(
		r.Context(),
		req.GUID.String(),
		sessionID,
		req.RefreshToken,
		userAgent,
		ip,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	guidVal := r.Context().Value(middleware.ContextUserGUIDKey)

	guid := guidVal.(string)

	response := MeResponse{
		GUID: guid,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
