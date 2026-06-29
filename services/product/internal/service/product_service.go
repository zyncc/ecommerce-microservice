package service

import (
	"context"

	"github.com/zyncc/ecommerce-microservice/services/product/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/product/internal/repository/model"
	"github.com/zyncc/ecommerce-microservice/services/product/pkg/types/dto"
	"go.uber.org/zap"
)

type ProductService struct {
	log  *zap.Logger
	repo *repository.ProductRepository
}

func NewProductService(log *zap.Logger, repo *repository.ProductRepository) *ProductService {
	return &ProductService{log, repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (string, error) {
	return s.repo.CreateProduct(ctx, &model.CreateProductParams{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
	})
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]*model.Product, error) {
	return s.repo.FetchAllProducts(ctx)
}
