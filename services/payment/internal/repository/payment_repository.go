package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zyncc/ecommerce-microservice/services/payment/internal/repository/models"
	"go.uber.org/zap"
)

type PaymentRepository struct {
	log *zap.Logger
	db  *pgxpool.Pool
}

func NewPaymentRepository(log *zap.Logger, db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{
		log,
		db,
	}
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, params *models.CreatePaymentParams) (uuid.UUID, error) {
	id := uuid.New()
	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO payments
		(
			id,
			order_id,
			idempotency_key,
			status,
			amount,
			payment_method,
			currency
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		`,
		id,
		params.OrderID,
		params.IdempotencyKey,
		params.Status,
		params.Amount,
		params.PaymentMethod,
		params.Currency,
	)
	if err != nil {
		r.log.Error("failed to create payment", zap.Error(err))
		return uuid.Nil, err
	}

	return id, nil
}

func (r *PaymentRepository) FindByIdempotencyKey(ctx context.Context, key uuid.UUID) (bool, error) {
	var exists bool

	err := r.db.QueryRow(
		ctx,
		`
		SELECT EXISTS (
			SELECT 1
			FROM payments
			WHERE idempotency_key = $1
		)
		`,
		key,
	).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
