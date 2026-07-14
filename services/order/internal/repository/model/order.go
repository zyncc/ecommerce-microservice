package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
)

type Order struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Subtotal       float64
	OrderTotal     float64
	ShippingCost   float64
	PaymentID      *uuid.UUID
	Waybill        *uuid.UUID
	OrderStatus    string
	PaymentStatus  string
	FirstName      string
	LastName       *string
	Email          string
	Phone          string
	Address1       string
	Address2       *string
	City           string
	State          string
	Zip            string
	IdempotencyKey *uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
}

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

type CreateOrderParams struct {
	UserID        uuid.UUID
	Items         []dto.OrderItem
	Subtotal      float64
	OrderTotal    float64
	ShippingCost  float64
	OrderStatus   string
	PaymentStatus string
	FirstName     string
	LastName      *string
	Email         string
	Phone         string
	Address1      string
	Address2      *string
	City          string
	State         string
	Zip           string
}
