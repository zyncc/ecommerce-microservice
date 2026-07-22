package models

import (
	"time"

	"github.com/google/uuid"
)

type Shipment struct {
	ID             uuid.UUID
	OrderID        uuid.UUID
	IdempotencyKey *uuid.UUID
	Status         string
	Carrier        string
	ShippingCost   float64
	TrackingNumber uuid.UUID

	ShippedAt   *time.Time
	DeliveredAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateShipmentParams struct {
	ID             uuid.UUID
	OrderID        uuid.UUID
	Carrier        string
	ShippingCost   float64
	TrackingNumber uuid.UUID
}

type UpdateShipmentParams struct {
	TrackingID     uuid.UUID
	IdempotencyKey uuid.UUID
	Status         string
	ShippedAt      *time.Time
	DeliveredAt    *time.Time
}
