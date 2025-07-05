package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/nerfthisdev/go-backend-test-task/internal/auth"
	"github.com/nerfthisdev/go-backend-test-task/internal/middleware"
)

// RefreshRequest represents a request payload for token refresh.
type RefreshRequest struct {
	GUID         uuid.UUID `json:"guid"`
	RefreshToken string    `json:"refresh_token"`
}

// TokenResponse contains generated JWT tokens.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// MeResponse represents the /me endpoint response.
type MeResponse struct {
	GUID uuid.UUID `json:"guid"`
}

type AuthHandler struct {
	auth *auth.AuthService
}

func NewAuthHandler(service *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		auth: service,
	}
}

// Authorize godoc
// @Summary      Authorize user
// @Description  Returns new access and refresh tokens. If guid is empty a new user is created.
// @Tags         auth
// @Param        guid  query     string  false  "User GUID"
// @Produce      json
// @Success      200  {object}  TokenResponse
// @Failure      401  {string}  string  "unauthorized"
// @Router       /auth [get]
func (h *AuthHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	guidStr := r.URL.Query().Get("guid")
	var guidPtr *uuid.UUID
	if guidStr != "" {
		parsed, err := uuid.Parse(guidStr)
		if err != nil {
			http.Error(w, "invalid guid", http.StatusBadRequest)
			return
		}
		guidPtr = &parsed
	}

	userAgent := r.UserAgent()
	ip := r.RemoteAddr

	tokens, err := h.auth.Authorize(r.Context(), guidPtr, userAgent, ip)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// Refresh godoc
// @Summary      Refresh tokens
// @Description  Generates new token pair using refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RefreshRequest true "Refresh request"
// @Success      200 {object} TokenResponse
// @Failure      400 {string} string "invalid request"
// @Failure      401 {string} string "unauthorized"
// @Security     BearerAuth
// @Router       /refresh [post]
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
		req.GUID,
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

// Me godoc
// @Summary      Get user info
// @Description  Returns current user GUID
// @Tags         auth
// @Produce      json
// @Success      200 {object} MeResponse
// @Failure      401 {string} string "unauthorized"
// @Security     BearerAuth
// @Router       /me [post]
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	guidVal := r.Context().Value(middleware.ContextUserGUIDKey)

	guid := guidVal.(uuid.UUID)

	response := MeResponse{
		GUID: guid,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) Deauthorize(w http.ResponseWriter, r *http.Request) {
	guidVal := r.Context().Value(middleware.ContextUserGUIDKey)

	guid := guidVal.(uuid.UUID)

	err := h.auth.Deauthorize(r.Context(), guid)
	if err != nil {
		http.Error(w, "failed to deauthorize", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
