package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	ProductID uuid.UUID
	Quantity  int
	Size      string
	Price     float64

	CreatedAt time.Time
	UpdatedAt time.Time
}
