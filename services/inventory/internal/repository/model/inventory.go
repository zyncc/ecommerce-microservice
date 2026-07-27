package model

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID         uuid.UUID
	ProductID  uuid.UUID
	Small      int
	Medium     int
	Large      int
	ExtraLarge int

	CreatedAt time.Time
	UpdatedAt time.Time
}
