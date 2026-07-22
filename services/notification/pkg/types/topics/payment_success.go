package topics

import (
	"time"

	"github.com/google/uuid"
)

type PaymentSucceededEvent struct {
	EventID       uuid.UUID `json:"event_id"`
	OrderID       uuid.UUID `json:"order_id"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	Currency      string    `json:"currency"`
	OccurredAt    time.Time `json:"occurred_at"`
}
