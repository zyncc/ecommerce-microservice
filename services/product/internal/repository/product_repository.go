package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zyncc/ecommerce-microservice/services/product/internal/repository/model"
	"github.com/zyncc/ecommerce-microservice/services/product/pkg/types"
	"go.uber.org/zap"
)

type ProductRepository struct {
	log *zap.Logger
	db  *pgxpool.Pool
}

func NewProductRepository(log *zap.Logger, db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{log, db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, params *model.CreateProductParams) (string, error) {
	id := uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO product (id, title, description, price, category)
		VALUES ($1, $2, $3, $4, $5)`,
		id, params.Title, params.Description, params.Price, params.Category)
	if err != nil {
		r.log.Error("failed to create product", zap.Error(err))
		return "", types.ErrDatabase
	}
	return id.String(), nil
}

func (r *ProductRepository) FetchAllProducts(ctx context.Context, limit, offset int) ([]*model.Product, error) {
	products := make([]*model.Product, 0)
	rows, err := r.db.Query(ctx,
		`SELECT id, title, description, price, category, created_at, updated_at
		FROM product
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		r.log.Error("failed to fetch products", zap.Error(err))
		return nil, types.ErrDatabase
	}
	defer rows.Close()

	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Title, &product.Description, &product.Price, &product.Category, &product.CreatedAt, &product.UpdatedAt); err != nil {
			r.log.Error("failed to scan product", zap.Error(err))
			return nil, types.ErrDatabase
		}
		products = append(products, &product)
	}
	if err := rows.Err(); err != nil {
		r.log.Error("error iterating product rows", zap.Error(err))
		return nil, types.ErrDatabase
	}
	return products, nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id uuid.UUID) (model.Product, error) {
	var product model.Product
	if err := r.db.QueryRow(ctx,
		`SELECT id, title, description, price, category, created_at, updated_at
		FROM product
		WHERE id = $1`,
		id,
	).Scan(&product.ID, &product.Title, &product.Description, &product.Price, &product.Category, &product.CreatedAt, &product.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Product{}, types.ErrProductNotFound
		}
		r.log.Error("failed to fetch product", zap.Error(err))
		return model.Product{}, types.ErrDatabase
	}

	return product, nil
}
