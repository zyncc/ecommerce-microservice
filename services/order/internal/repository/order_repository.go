package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zyncc/ecommerce-microservice/services/order/internal/repository/model"
	"github.com/zyncc/ecommerce-microservice/services/order/pkg/types"
	"go.uber.org/zap"
)

type OrderRepository struct {
	log *zap.Logger
	db  *pgxpool.Pool
}

func NewOrderRepository(log *zap.Logger, db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{log, db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, params model.CreateOrderParams) (uuid.UUID, error) {
	orderID := uuid.New()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		r.log.Error("failed to create transaction", zap.Error(err))
		return uuid.Nil, types.ErrDatabase
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		`INSERT INTO orders (
			id,
			user_id,
			subtotal,
			order_total,
			shipping_cost,
			first_name,
			last_name,
			email,
			phone,
			address1,
			address2,
			city,
			state,
			zip
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13, $14
		)`,
		orderID,
		params.UserID,
		params.Subtotal,
		params.OrderTotal,
		params.ShippingCost,
		params.FirstName,
		params.LastName,
		params.Email,
		params.Phone,
		params.Address1,
		params.Address2,
		params.City,
		params.State,
		params.Zip,
	)
	if err != nil {
		r.log.Error("failed to create order", zap.Error(err))
		return uuid.Nil, types.ErrDatabase
	}

	for _, item := range params.Items {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO order_items (
				id,
				order_id,
				product_id,
				quantity,
				size,
				price
			) VALUES ($1, $2, $3, $4, $5, $6)`,
			uuid.New(),
			orderID,
			item.ProductID,
			item.Quantity,
			item.Size,
			item.Price,
		)
		if err != nil {
			r.log.Error("failed to create order_items", zap.Error(err))
			return uuid.Nil, types.ErrDatabase
		}
	}

	if err := tx.Commit(ctx); err != nil {
		r.log.Error("failed to commit transaction", zap.Error(err))
		return uuid.Nil, types.ErrDatabase
	}

	return orderID, nil
}
