package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID             uuid.UUID
	OrderID        uuid.UUID
	IdempotencyKey uuid.UUID
	Status         string
	Amount         float64
	PaymentMethod  string
	Currency       string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreatePaymentParams struct {
	OrderID        uuid.UUID
	Status         string
	Amount         float64
	PaymentMethod  string
	Currency       string
	IdempotencyKey uuid.UUID
}
