package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/nerfthisdev/go-backend-test-task/internal/auth"
)

type AuthHandler struct {
	service *auth.AuthService
}

func NewAuthHandler(service *auth.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	userID := uuid.New()

	tokenPair, err := h.service.CreateTokenPair(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error": "failed to create token pair"`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tokenPair); err != nil {
		http.Error(w, `{"error": "failed to write response"}`, http.StatusInternalServerError)
	}
}

func (h *AuthHandler) FetchUser(w http.ResponseWriter, r *http.Request) {
	guidString := r.PathValue("id")
	guid, err := uuid.Parse(guidString)
	if err != nil {
		http.Error(w, `{"error": "malformed guid"}`, http.StatusBadRequest)
	}

	userResponse, err := h.service.GetUserByGuid(r.Context(), guid)
	if err != nil {
		http.Error(w, `"error": "failed to fetch user"`, http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userResponse); err != nil {
		http.Error(w, `{"error": "failed to write response"}`, http.StatusInternalServerError)
	}
}
