package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
)

type Order struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Subtotal   float64
	OrderTotal float64

	OrderStatus string
	FirstName   string
	LastName    *string
	Email       string
	Phone       string
	Address1    string
	Address2    *string
	City        string
	State       string
	Zip         string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateOrderParams struct {
	UserID     uuid.UUID
	Items      []dto.OrderItem
	Subtotal   float64
	OrderTotal float64
	FirstName  string
	LastName   *string
	Email      string
	Phone      string
	Address1   string
	Address2   *string
	City       string
	State      string
	Zip        string
}

type OrderWithItems struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Subtotal    float64
	OrderTotal  float64
	OrderStatus string
	FirstName   string
	LastName    *string
	Email       string
	Phone       string
	Address1    string
	Address2    *string
	City        string
	State       string
	Zip         string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	OrderItems []OrderItem
}
