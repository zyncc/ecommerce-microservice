package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	pb "github.com/zyncc/ecommerce-microservice/services/product/pkg/types/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductClient interface {
	CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (uuid.UUID, error)
	GetAllProducts(ctx context.Context, limit, offset int) ([]dto.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (dto.Product, error)
}

type GRPCProductClient struct {
	log           *zap.Logger
	productSvcURL string
	client        pb.ProductServiceClient
}

func NewGRPCProductClient(log *zap.Logger, addr string) (*GRPCProductClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GRPCProductClient{
		log:           log,
		productSvcURL: addr,
		client:        pb.NewProductServiceClient(conn),
	}, nil
}

func (c *GRPCProductClient) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (uuid.UUID, error) {
	resp, err := c.client.CreateProduct(ctx, &pb.CreateProductRequest{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Inventory: &pb.Inventory{
			Small:      int32(*req.Inventory.Small),
			Medium:     int32(*req.Inventory.Medium),
			Large:      int32(*req.Inventory.Large),
			ExtraLarge: int32(*req.Inventory.ExtraLarge),
		},
	})
	if err != nil {
		return uuid.Nil, err
	}

	productID, _ := uuid.Parse(resp.GetId())

	return productID, nil
}

func (c *GRPCProductClient) GetAllProducts(ctx context.Context, limit, offset int) ([]dto.Product, error) {
	resp, err := c.client.GetAllProducts(ctx, &pb.GetAllProductsRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	var products []dto.Product
	for _, item := range resp.GetProducts() {
		productID, _ := uuid.Parse(item.GetId())
		product := dto.Product{
			ID:          productID,
			Title:       item.GetTitle(),
			Description: item.GetDescription(),
			Price:       item.GetPrice(),
			Category:    item.GetCategory(),
			CreatedAt:   item.CreatedAt.AsTime(),
			UpdatedAt:   item.UpdatedAt.AsTime(),
		}

		products = append(products, product)
	}

	return products, nil
}

func (c *GRPCProductClient) GetProductByID(ctx context.Context, id uuid.UUID) (dto.Product, error) {
	resp, err := c.client.GetProductByID(ctx, &pb.IDMessage{Id: id.String()})
	if err != nil {
		return dto.Product{}, err
	}

	productID, _ := uuid.Parse(resp.GetId())
	product := dto.Product{
		ID:          productID,
		Title:       resp.GetTitle(),
		Description: resp.GetDescription(),
		Price:       resp.GetPrice(),
		Category:    resp.GetCategory(),
		CreatedAt:   resp.GetCreatedAt().AsTime(),
		UpdatedAt:   resp.GetUpdatedAt().AsTime(),
	}

	return product, nil
}
