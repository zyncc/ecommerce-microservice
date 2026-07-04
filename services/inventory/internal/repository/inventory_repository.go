package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/repository/model"
	"github.com/zyncc/ecommerce-microservice/services/inventory/pkg/types"
	"go.uber.org/zap"
)

type InventoryRepository struct {
	log *zap.Logger
	db  *pgxpool.Pool
}

func NewInventoryRepository(log *zap.Logger, db *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{log, db}
}

func (r *InventoryRepository) CreateInventory(ctx context.Context, params *model.CreateInventoryParams) (uuid.UUID, error) {
	id := uuid.New()

	_, err := r.db.Exec(ctx,
		`INSERT INTO inventory
		(id, product_id, small, medium, large, extra_large)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		id, params.ProductID, params.Small, params.Medium, params.Large, params.ExtraLarge)
	if err != nil {
		r.log.Error("failed to create inventory", zap.Error(err))
		return uuid.Nil, types.ErrDatabase
	}

	return id, nil
}

func (r *InventoryRepository) FindInventoryByProductID(ctx context.Context, productID uuid.UUID) (model.Inventory, error) {
	var inventory model.Inventory

	if err := r.db.QueryRow(ctx,
		`SELECT id, product_id, small, medium, large, extra_large, created_at, updated_at
		FROM inventory WHERE product_id = $1`,
		productID,
	).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.Small,
		&inventory.Medium,
		&inventory.Large,
		&inventory.ExtraLarge,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Inventory{}, types.ErrInventoryNotFound
		}
		return model.Inventory{}, types.ErrDatabase
	}

	return inventory, nil
}

func (r *InventoryRepository) UpdateQuantity(
	ctx context.Context,
	productID uuid.UUID,
	size string,
	quantity int,
) error {
	switch size {
	case "small", "medium", "large", "extra_large":
	default:
		return types.ErrInvalidSize
	}

	query := fmt.Sprintf(`
		UPDATE inventory
		SET %s = $1,
		updated_at = NOW()
		WHERE product_id = $2`,
		size)

	tag, err := r.db.Exec(ctx, query, quantity, productID)
	if err != nil {
		return types.ErrDatabase
	}

	if tag.RowsAffected() == 0 {
		return types.ErrInventoryNotFound
	}

	return nil
}
