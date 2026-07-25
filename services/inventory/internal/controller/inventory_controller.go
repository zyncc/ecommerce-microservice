package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/service"
	pb "github.com/zyncc/ecommerce-microservice/services/inventory/pkg/types/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InventoryController struct {
	pb.UnimplementedInventoryServiceServer
	log *zap.Logger
	svc *service.InventoryService
}

func NewInventoryController(log *zap.Logger, svc *service.InventoryService) *InventoryController {
	return &InventoryController{
		log: log,
		svc: svc,
	}
}

func (c *InventoryController) GetInventoryByProductID(ctx context.Context, in *pb.IDMessage) (*pb.GetInventoryResponse, error) {
	productID, _ := uuid.Parse(in.GetId())

	inventory, err := c.svc.FetchInventoryByProductID(ctx, productID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetInventoryResponse{
		Id:         inventory.ID.String(),
		ProductId:  inventory.ProductID.String(),
		Small:      int32(inventory.Small),
		Medium:     int32(inventory.Medium),
		Large:      int32(inventory.Large),
		ExtraLarge: int32(inventory.ExtraLarge),
		CreatedAt:  timestamppb.New(inventory.CreatedAt),
		UpdatedAt:  timestamppb.New(inventory.UpdatedAt),
	}, nil
}
