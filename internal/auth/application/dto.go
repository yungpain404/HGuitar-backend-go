package application

import (
	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/auth/infrastructure/security"
)

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	DisplayName string `json:"display_name"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type FirebaseLoginRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

type AuthResponse struct {
	UserID      string   `json:"user_id"`
	Email       string   `json:"email"`
	DisplayName string   `json:"display_name"`
	Roles       []string `json:"roles"`
	AccessToken string   `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
}

type FirebaseIdentity = security.FirebaseIdentity

type AuthService interface {
	Register(req RegisterRequest) (*AuthResponse, error)
	Login(req LoginRequest) (*AuthResponse, error)
	LoginWithFirebase(identity *FirebaseIdentity) (*AuthResponse, error)
	Refresh(req RefreshRequest) (*AuthResponse, error)
	Logout(req LogoutRequest) error
}

type UserProfileService interface {
	CreateProfile(userID uuid.UUID, displayName, avatarURL string) error
	GetDisplayName(userID uuid.UUID) (string, error)
}
