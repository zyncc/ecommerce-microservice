package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	Title       string
	Description string
	Price       float64
	Category    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateProductParams struct {
	Title       string
	Description string
	Price       float64
	Category    string
}
