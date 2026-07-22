package dto

import (
	"time"

	"github.com/google/uuid"
)

type ShipmentWebhookRequest struct {
	IdempotencyKey uuid.UUID `json:"idempotency_key"`
	TrackingNumber uuid.UUID `json:"tracking_number"`
	Status         string    `json:"status"`
	OccuredAt      time.Time `json:"occured_at"`
}

type ShipmentResponse struct {
	ID             uuid.UUID  `json:"id"`
	OrderID        uuid.UUID  `json:"order_id"`
	IdempotencyKey *uuid.UUID `json:"idempotency_key"`
	Status         string     `json:"status"`
	Carrier        string     `json:"carrier"`
	ShippingCost   float64    `json:"shipping_cost"`
	TrackingNumber uuid.UUID  `json:"tracking_number"`

	ShippedAt   *time.Time `json:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
