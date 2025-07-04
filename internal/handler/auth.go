package handler

import "github.com/google/uuid"

type RefreshRequest struct {
	GUID         uuid.UUID `json:"guid"`
	RefreshToken string    `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
