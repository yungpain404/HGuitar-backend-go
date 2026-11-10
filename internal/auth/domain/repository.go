package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	ExistsByEmail(email string) (bool, error)
	FindByEmail(email string) (*User, error)
	FindByID(id uuid.UUID) (*User, error)
	Save(user *User) error
}

type RefreshTokenRepository interface {
	FindUserIDByTokenHash(tokenHash string) (uuid.UUID, error)
	RevokeByTokenHash(tokenHash string) error
	Save(userID uuid.UUID, tokenHash string, expiresAt time.Time) error
}
