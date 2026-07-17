package repository

import (
	"context"
	"errors"

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

var ErrOrderNotFound = errors.New("order not found")

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
			$8, $9, $10, $11, $12, $13
		)`,
		orderID,
		params.UserID,
		params.Subtotal,
		params.OrderTotal,
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

func (r *OrderRepository) GetOrder(ctx context.Context, orderID uuid.UUID) (model.OrderWithItems, error) {
	var order model.OrderWithItems
	query := `
		SELECT
		o.id,
		o.user_id,
		o.idempotency_key,
		o.subtotal,
		o.order_total,
		o.order_status,
		o.first_name,
		o.last_name,
		o.email,
		o.phone,
		o.address1,
		o.address2,
		o.city,
		o.state,
		o.zip,
		o.created_at,
		o.updated_at,

		oi.id,
		oi.order_id,
    oi.product_id,
    oi.quantity,
    oi.size,
    oi.price,
		oi.created_at,
		oi.updated_at

		FROM orders AS o
		INNER JOIN order_items AS oi
		ON oi.order_id = o.id
		WHERE o.id = $1
		`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		r.log.Error("failed to get order", zap.Error(err))
		return model.OrderWithItems{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem

		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.IdempotencyKey,
			&order.Subtotal,
			&order.OrderTotal,
			&order.OrderStatus,
			&order.FirstName,
			&order.LastName,
			&order.Email,
			&order.Phone,
			&order.Address1,
			&order.Address2,
			&order.City,
			&order.State,
			&order.Zip,
			&order.CreatedAt,
			&order.UpdatedAt,

			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Size,
			&item.Price,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			r.log.Error("failed to scan order into struct", zap.Error(err))
			return model.OrderWithItems{}, err
		}

		order.OrderItems = append(order.OrderItems, item)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("failed to iterate rows", zap.Error(err))
		return model.OrderWithItems{}, err
	}

	if order.ID == uuid.Nil {
		return model.OrderWithItems{}, ErrOrderNotFound
	}

	return order, nil
}

func (r *OrderRepository) UpdateIdempotencyKeyAndOrderStatus(ctx context.Context, orderID uuid.UUID, key uuid.UUID, status string) error {
	tag, err := r.db.Exec(
		ctx,
		`
		UPDATE orders
		SET 
		idempotency_key = $1,
		order_status = $2
		WHERE id = $3
		`,
		key, status, orderID,
	)
	if err != nil {
		r.log.Error("failed to update idempotency_key and order_status", zap.Error(err))
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("order not found")
	}

	return nil
}
