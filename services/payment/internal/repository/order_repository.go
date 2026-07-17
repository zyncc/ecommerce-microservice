package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type OrderRepository struct {
	log *zap.Logger
	db  *pgxpool.Pool
}

func NewOrderRepository(log *zap.Logger, db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		log,
		db,
	}
}

func (r *OrderRepository) FindOrderByIdempotencyKey(ctx context.Context, key uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		ctx,
		`
		SELECT EXISTS (
			SELECT 1
			FROM orders
			WHERE idempotency_key = $1
		)
		`,
		key,
	).Scan(&exists)
	if err != nil {
		r.log.Error("failed to find order by idempotency_key", zap.Error(err))
		return false, errors.New("failed to find order by idempotency_key")
	}

	return exists, nil
}
