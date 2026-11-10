package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	DisplayName string
	AvatarURL   string
	GuitarType  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewUserProfile(id, userID uuid.UUID, displayName, avatarURL, guitarType string, now time.Time) *UserProfile {
	return &UserProfile{
		ID:          id,
		UserID:      userID,
		DisplayName: displayName,
		AvatarURL:   avatarURL,
		GuitarType:  guitarType,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
