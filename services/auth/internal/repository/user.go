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

type UserRepository struct {
	logger *zap.Logger
	db     *pgxpool.Pool
}

func NewUserRepository(logger *zap.Logger, db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		logger,
		db,
	}
}

var ErrUserNotFound = errors.New("user not found")

func (r *UserRepository) CreateUser(ctx context.Context, params *models.CreateUserParams) error {
	_, err := r.db.Exec(ctx, "INSERT INTO users (id, name, email, hashed_password, role) VALUES ($1, $2, $3, $4, $5)", params.ID, params.Name, params.Email, params.HashedPassword, params.Role)
	if err != nil {
		r.logger.Error("failed to create user", zap.Error(err))
		return err
	}

	return nil
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		ctx,
		"SELECT id, name, email, hashed_password, role, created_at, updated_at FROM users WHERE email = $1",
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Error("user with this email does not exist", zap.String("email", email), zap.Error(err))
		return nil, ErrUserNotFound
	} else if err != nil {
		r.logger.Error("failed to fetch user", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		ctx,
		"SELECT id, name, email, hashed_password, role, created_at, updated_at FROM users WHERE id = $1",
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Error("user with this id does not exist", zap.String("id", id.String()), zap.Error(err))
		return nil, ErrUserNotFound
	} else if err != nil {
		r.logger.Error("failed to fetch user", zap.Error(err))
		return nil, err
	}

	return &user, nil
}
