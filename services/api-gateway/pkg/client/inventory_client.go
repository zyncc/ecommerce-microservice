package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	pb "github.com/zyncc/ecommerce-microservice/services/inventory/pkg/types/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type InventoryClient interface {
	FetchInventoryByProductID(ctx context.Context, productID uuid.UUID) (dto.InventoryResponse, error)
}

type GRPCInventoryClient struct {
	log             *zap.Logger
	inventorySvcURL string
	client          pb.InventoryServiceClient
}

func NewGRPCInventoryClient(log *zap.Logger, addr string) (*GRPCInventoryClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GRPCInventoryClient{
		log:             log,
		inventorySvcURL: addr,
		client:          pb.NewInventoryServiceClient(conn),
	}, nil
}

func (c *GRPCInventoryClient) FetchInventoryByProductID(ctx context.Context, productID uuid.UUID) (dto.InventoryResponse, error) {
	resp, err := c.client.GetInventoryByProductID(ctx, &pb.IDMessage{Id: productID.String()})
	if err != nil {
		return dto.InventoryResponse{}, err
	}

	ID, _ := uuid.Parse(resp.GetId())
	productID, _ = uuid.Parse(resp.GetProductId())
	dto := dto.InventoryResponse{
		ID:         ID,
		ProductID:  productID,
		Small:      int(resp.GetSmall()),
		Medium:     int(resp.GetMedium()),
		Large:      int(resp.GetLarge()),
		ExtraLarge: int(resp.GetExtraLarge()),
		CreatedAt:  resp.GetCreatedAt().AsTime(),
		UpdatedAt:  resp.GetUpdatedAt().AsTime(),
	}

	return dto, nil
}
