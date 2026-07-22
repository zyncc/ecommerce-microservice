package types

import (
	"time"

	"github.com/google/uuid"
)

const (
	PaymentSucceededTopic = "payment.succeeded"
	ShipmentUpdatedTopic  = "shipment.updated"
)

type ShipmentWebhookRequest struct {
	IdempotencyKey uuid.UUID `json:"idempotency_key"`
	TrackingNumber uuid.UUID `json:"tracking_number"`
	Status         string    `json:"status"`
	OccuredAt      time.Time `json:"occured_at"`
}
