package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	Name           string
	Email          string
	HashedPassword string
	Role           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreateUserParams struct {
	ID             uuid.UUID
	Name           string
	Email          string
	HashedPassword string
	Role           string
}

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}
