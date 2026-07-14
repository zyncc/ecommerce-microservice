package dto

import "github.com/google/uuid"

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
