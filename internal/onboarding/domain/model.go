package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserPreference struct {
	ID   uuid.UUID
	UserID uuid.UUID
	GuitarType string
	Completed bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserPreferenceRepository interface {
	FindByUserID(userID uuid.UUID) (*UserPreference, error)
	Save(pref *UserPreference) error
}
