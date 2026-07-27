package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/product/internal/service"
	pb "github.com/zyncc/ecommerce-microservice/services/product/pkg/types/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductController struct {
	pb.UnimplementedProductServiceServer
	log *zap.Logger
	svc *service.ProductService
}

func NewProductController(log *zap.Logger, svc *service.ProductService) *ProductController {
	return &ProductController{
		log: log,
		svc: svc,
	}
}

func (c *ProductController) CreateProduct(ctx context.Context, in *pb.CreateProductRequest) (*pb.IDMessage, error) {
	small := int(in.GetInventory().GetSmall())
	medium := int(in.GetInventory().GetMedium())
	large := int(in.GetInventory().GetLarge())
	extraLarge := int(in.GetInventory().GetExtraLarge())

	id, err := c.svc.CreateProduct(ctx, &dto.CreateProductRequest{
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
		Price:       in.GetPrice(),
		Category:    in.GetCategory(),
		Inventory: &dto.Inventory{
			Small:      &small,
			Medium:     &medium,
			Large:      &large,
			ExtraLarge: &extraLarge,
		},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create product")
	}

	return &pb.IDMessage{
		Id: id.String(),
	}, nil
}

func (c *ProductController) GetAllProducts(ctx context.Context, in *pb.GetAllProductsRequest) (*pb.GetAllProductsResponse, error) {
	const (
		defaultLimit = 5
		maxLimit     = 10
	)

	limit := int(in.GetLimit())

	switch {
	case limit == 0:
		limit = defaultLimit
	case limit <= 0:
		return nil, status.Error(codes.InvalidArgument, "limit must be greater than 0")
	case limit > maxLimit:
		limit = maxLimit
	}

	offset := int(in.GetOffset())
	if offset < 0 {
		return nil, status.Error(codes.InvalidArgument, "offset must be greater than or equal to 0")
	}

	productsResp, err := c.svc.GetAllProducts(ctx, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to fetch all products")
	}

	var products []*pb.Product
	for _, p := range productsResp {
		product := &pb.Product{
			Id:          p.ID.String(),
			Title:       p.Title,
			Description: p.Description,
			Price:       p.Price,
			Category:    p.Category,
			CreatedAt:   timestamppb.New(p.CreatedAt),
			UpdatedAt:   timestamppb.New(p.UpdatedAt),
		}

		products = append(products, product)
	}

	return &pb.GetAllProductsResponse{
		Products: products,
	}, nil
}

func (c *ProductController) GetProductByID(ctx context.Context, in *pb.IDMessage) (*pb.Product, error) {
	productID, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "product id must be a valid uuid")
	}

	product, err := c.svc.GetProductByID(ctx, productID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get product by id")
	}

	return &pb.Product{
		Id:          product.ID.String(),
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
	}, nil
}
