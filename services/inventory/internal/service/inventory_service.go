package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/repository/model"
	"github.com/zyncc/ecommerce-microservice/services/inventory/pkg/types/dto"
	"go.uber.org/zap"
)

type InventoryService struct {
	log  *zap.Logger
	repo *repository.InventoryRepository
}

func NewInventoryService(log *zap.Logger, repo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{log, repo}
}

func (s *InventoryService) CreateInventory(ctx context.Context, req *dto.CreateInventoryRequest) (uuid.UUID, error) {
	return s.repo.CreateInventory(ctx, &model.CreateInventoryParams{
		ProductID:  req.ProductID,
		Small:      req.Inventory.Small,
		Medium:     req.Inventory.Medium,
		Large:      req.Inventory.Large,
		ExtraLarge: req.Inventory.ExtraLarge,
	})
}

func (s *InventoryService) FetchInventoryByProductID(ctx context.Context, productID uuid.UUID) (dto.InventoryResponse, error) {
	inventory, err := s.repo.FindInventoryByProductID(ctx, productID)
	if err != nil {
		return dto.InventoryResponse{}, err
	}

	response := dto.InventoryResponse{
		ID:         inventory.ID,
		ProductID:  inventory.ProductID,
		Small:      inventory.Small,
		Medium:     inventory.Medium,
		Large:      inventory.Large,
		ExtraLarge: inventory.ExtraLarge,
		CreatedAt:  inventory.CreatedAt,
		UpdatedAt:  inventory.UpdatedAt,
	}

	return response, nil
}

func (s *InventoryService) UpdateInventory(ctx context.Context, req []dto.UpdateInventoryRequest) error {
	return s.repo.UpdateInventory(ctx, req)
}
