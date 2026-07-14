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

type CreateInventoryRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Inventory Inventory `json:"inventory"`
}

type UpdateInventoryRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Size      string    `json:"size"`
	Quantity  int       `json:"quantity"`
}

type InventoryResponse struct {
	ID         uuid.UUID `json:"id"`
	ProductID  uuid.UUID `json:"product_id"`
	Small      int       `json:"small"`
	Medium     int       `json:"medium"`
	Large      int       `json:"large"`
	ExtraLarge int       `json:"extra_large"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
