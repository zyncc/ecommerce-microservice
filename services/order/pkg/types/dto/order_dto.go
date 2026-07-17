package dto

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	Small      int `json:"small"`
	Medium     int `json:"medium"`
	Large      int `json:"large"`
	ExtraLarge int `json:"extra_large"`
}

type UpdateInventoryRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Size      string    `json:"size"`
	Quantity  int       `json:"quantity"`
}

type FindOrderByIDResponse struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	IdempotencyKey *uuid.UUID `json:"idempotency_key"`
	Subtotal       float64    `json:"subtotal"`
	OrderTotal     float64    `json:"order_total"`
	OrderStatus    string     `json:"order_status"`
	FirstName      string     `json:"first_name"`
	LastName       *string    `json:"last_name"`
	Email          string     `json:"email"`
	Phone          string     `json:"phone"`
	Address1       string     `json:"address_1"`
	Address2       *string    `json:"address_2"`
	City           string     `json:"city"`
	State          string     `json:"state"`
	Zip            string     `json:"zip"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	OrderItems []OrderItems `json:"order_items"`
}

type OrderItems struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Size      string    `json:"size"`
	Price     float64   `json:"price"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
