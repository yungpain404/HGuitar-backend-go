package domain

import (
	"github.com/google/uuid"
)

type UserProfileRepository interface {
	FindByUserID(userID uuid.UUID) (*UserProfile, error)
	Save(profile *UserProfile) error
}
