package dto

import (
	"github.com/google/uuid"
)

type PaymentWebhookRequest struct {
	IdempotencyKey uuid.UUID `json:"idempotency_key"`
	OrderID        uuid.UUID `json:"order_id"`
	Amount         float64   `json:"amount"`
	PaymentMethod  string    `json:"payment_method"`
	Currency       string    `json:"currency"`
	Status         string    `json:"status"`
}
