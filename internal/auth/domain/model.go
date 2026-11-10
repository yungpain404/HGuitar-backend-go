package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
)

type Role = common.Role

const (
	RoleStudent       = common.RoleStudent
	RoleContentEditor = common.RoleContentEditor
	RoleAdmin         = common.RoleAdmin
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Active       bool
	Roles        []Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(id uuid.UUID, email, passwordHash string, active bool, roles []Role, now time.Time) *User {
	return &User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
		Active:       active,
		Roles:        roles,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (u *User) Deactivate() {
	u.Active = false
	u.UpdatedAt = time.Now()
}

type RefreshToken struct {
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
}
