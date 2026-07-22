package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zyncc/ecommerce-microservice/services/shipping/internal/repository/models"
	"go.uber.org/zap"
)

type ShipmentRepository struct {
	log *zap.Logger
	db  *pgxpool.Pool
}

func NewShippingRepository(log *zap.Logger, db *pgxpool.Pool) *ShipmentRepository {
	return &ShipmentRepository{
		log,
		db,
	}
}

func (r *ShipmentRepository) CreateShipment(ctx context.Context, params *models.CreateShipmentParams) (uuid.UUID, error) {
	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO shipments
		(
			id,
			order_id,
			carrier,
			shipping_cost,
			tracking_number
		)
		VALUES (
			$1, $2, $3, $4, $5
		)
		`,
		params.ID,
		params.OrderID,
		params.Carrier,
		params.ShippingCost,
		params.TrackingNumber,
	)
	if err != nil {
		r.log.Error("failed to create shipment", zap.Error(err))
		return uuid.Nil, err
	}

	return params.ID, nil
}

func (r *ShipmentRepository) FindShipmentByIdempotencyKey(ctx context.Context, idempotencyKey uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		ctx,
		`
		SELECT EXISTS (
			SELECT 1
			FROM shipments
			WHERE idempotency_key = $1
		)
		`,
		idempotencyKey.String(),
	).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, err
}

func (r *ShipmentRepository) GetShipmentByTrackingID(ctx context.Context, trackingID uuid.UUID) (models.Shipment, error) {
	var shipment models.Shipment
	err := r.db.QueryRow(
		ctx,
		`SELECT
		id,
		order_id,
		idempotency_key,
		status,
		carrier,
		shipping_cost,
		tracking_number,
		shipped_at,
		delivered_at,
		created_at,
		updated_at
		FROM shipments
		WHERE tracking_number = $1
		`,
		trackingID.String(),
	).Scan(
		&shipment.ID,
		&shipment.OrderID,
		&shipment.IdempotencyKey,
		&shipment.Status,
		&shipment.Carrier,
		&shipment.ShippingCost,
		&shipment.TrackingNumber,
		&shipment.ShippedAt,
		&shipment.DeliveredAt,
		&shipment.CreatedAt,
		&shipment.UpdatedAt,
	)
	if err != nil {
		r.log.Error("failed to fetch shipment by tracking_number", zap.String("tracking_number", trackingID.String()), zap.Error(err))
		return models.Shipment{}, err
	}

	return shipment, nil
}

func (r *ShipmentRepository) UpdateShipment(ctx context.Context, params models.UpdateShipmentParams) error {
	query := `
		UPDATE shipments
		SET 
		status = $1,
    shipped_at = COALESCE($2, shipped_at),
		delivered_at = COALESCE($3, delivered_at),
		idempotency_key = $4,
		updated_at = NOW()
		WHERE tracking_number = $5
		`

	tag, err := r.db.Exec(ctx, query, params.Status, params.ShippedAt, params.DeliveredAt, params.IdempotencyKey, params.TrackingID)
	if err != nil {
		r.log.Error("failed to update shipment", zap.Error(err))
		return err
	}

	if tag.RowsAffected() == 0 {
		r.log.Error("shipment not found")
		return err
	}

	return nil
}
