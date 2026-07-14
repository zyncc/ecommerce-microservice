package dto

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Category    string     `json:"category"`
	Inventory   *Inventory `json:"inventory"`
}

func (r *CreateProductRequest) Validate() map[string][]string {
	errs := make(map[string][]string)

	if strings.TrimSpace(r.Title) == "" {
		errs["title"] = append(errs["title"], "title is required")
	}

	if strings.TrimSpace(r.Description) == "" {
		errs["description"] = append(errs["description"], "description is required")
	}

	if r.Price <= 0 {
		errs["price"] = append(errs["price"], "price must be greater than 0")
	}

	if strings.TrimSpace(r.Category) == "" {
		errs["category"] = append(errs["category"], "category is required")
	}

	if r.Inventory == nil {
		errs["inventory"] = append(errs["inventory"], "inventory is required")
	} else {
		if r.Inventory.Small == nil {
			errs["small"] = append(errs["small"], "small quantity is required")
		} else if *r.Inventory.Small < 0 {
			errs["small"] = append(errs["small"], "quantity cannot be below 0")
		}

		if r.Inventory.Medium == nil {
			errs["medium"] = append(errs["medium"], "medium quantity is required")
		} else if *r.Inventory.Medium < 0 {
			errs["medium"] = append(errs["medium"], "quantity cannot be below 0")
		}

		if r.Inventory.Large == nil {
			errs["large"] = append(errs["large"], "large quantity is required")
		} else if *r.Inventory.Large < 0 {
			errs["large"] = append(errs["large"], "quantity cannot be below 0")
		}

		if r.Inventory.ExtraLarge == nil {
			errs["extra_large"] = append(errs["extra_large"], "extra large quantity is required")
		} else if *r.Inventory.ExtraLarge < 0 {
			errs["extra_large"] = append(errs["extra_large"], "quantity cannot be below 0")
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
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
