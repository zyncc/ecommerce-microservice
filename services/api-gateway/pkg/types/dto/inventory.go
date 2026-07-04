package dto

import "github.com/google/uuid"

type CreateInventoryRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Inventory Inventory `json:"inventory"`
}
