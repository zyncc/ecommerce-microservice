package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zyncc/ecommerce-microservice/services/product/internal/repository/model"
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
	_, err := r.db.Exec(ctx, "INSERT INTO product (id, title, description, price, category) VALUES ($1, $2, $3, $4, $5)", id, params.Title, params.Description, params.Price, params.Category)
	if err != nil {
		r.log.Error("failed to create product", zap.Error(err))
		return "", errors.New("failed to create product")
	}
	return id.String(), nil
}

func (r *ProductRepository) FetchAllProducts(ctx context.Context) ([]*model.Product, error) {
	var products []*model.Product
	rows, err := r.db.Query(ctx, "SELECT id, title, description, price, category, created_at, updated_at FROM product")
	if err != nil {
		r.log.Error("failed to fetch products", zap.Error(err))
		return nil, errors.New("failed to fetch products")
	}

	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Title, &product.Description, &product.Price, &product.Category, &product.CreatedAt, &product.UpdatedAt); err != nil {
			r.log.Error("failed to scan product", zap.Error(err))
			return nil, errors.New("failed to fetch products")
		}
		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("error iterating product rows", zap.Error(err))
		return nil, errors.New("failed to fetch products")
	}

	return products, nil
}
