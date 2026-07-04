package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	Inventory   Inventory `json:"inventory"`
}

type Inventory struct {
	Small      int `json:"small"`
	Medium     int `json:"medium"`
	Large      int `json:"large"`
	ExtraLarge int `json:"extra_large"`
}

type Product struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
