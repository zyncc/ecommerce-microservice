package dto

import "github.com/google/uuid"

type CreateOrderRequest struct {
	Items     []OrderItem `json:"items"`
	UserID    uuid.UUID   `json:"user_id"`
	AddressID uuid.UUID   `json:"address_id"`
}

type OrderItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Size      string    `json:"size"`
	Price     float64
}
