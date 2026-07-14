package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/repository/models"
	"go.uber.org/zap"
)

type AddressRepository struct {
	logger *zap.Logger
	db     *pgxpool.Pool
}

func NewAddressRepository(logger *zap.Logger, db *pgxpool.Pool) *AddressRepository {
	return &AddressRepository{
		logger,
		db,
	}
}

var ErrAddressNotFound = errors.New("address not found")

func (r *AddressRepository) CreateAddress(ctx context.Context, params *models.CreateAddressParams) error {
	_, err := r.db.Exec(
		ctx, `
		INSERT INTO address (
			id,
			user_id,
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
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11
		)
	`,
		params.ID,
		params.UserID,
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
		r.logger.Error("failed to create address", zap.Error(err))
		return err
	}

	return nil
}

func (r *AddressRepository) FindAddressByID(ctx context.Context, id uuid.UUID) (*models.Address, error) {
	var address models.Address
	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			user_id,
			first_name,
			last_name,
			email,
			phone,
			address1,
			address2,
			city,
			state,
			zip,
			created_at,
			updated_at
		FROM address
		WHERE id = $1
	`, id).Scan(
		&address.ID,
		&address.UserID,
		&address.FirstName,
		&address.LastName,
		&address.Email,
		&address.Phone,
		&address.Address1,
		&address.Address2,
		&address.City,
		&address.State,
		&address.Zip,
		&address.CreatedAt,
		&address.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Error("address not found", zap.String("id", id.String()), zap.Error(err))
			return nil, ErrAddressNotFound
		}
		r.logger.Error("failed to fetch address", zap.Error(err))
		return nil, err
	}

	return &address, nil
}

func (r *AddressRepository) FetchAllAddresses(ctx context.Context, userID uuid.UUID) ([]models.Address, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			user_id,
			first_name,
			last_name,
			email,
			phone,
			address1,
			address2,
			city,
			state,
			zip,
			created_at,
			updated_at
		FROM address
		WHERE user_id = $1
	`, userID)
	if err != nil {
		r.logger.Error("failed to fetch all addresses", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	addresses := make([]models.Address, 0)

	for rows.Next() {
		var address models.Address

		err := rows.Scan(
			&address.ID,
			&address.UserID,
			&address.FirstName,
			&address.LastName,
			&address.Email,
			&address.Phone,
			&address.Address1,
			&address.Address2,
			&address.City,
			&address.State,
			&address.Zip,
			&address.CreatedAt,
			&address.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("failed to scan rows into address", zap.Error(err))
			return nil, err
		}
		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("failed to fetch address", zap.Error(err))
		return nil, err
	}

	return addresses, nil
}
