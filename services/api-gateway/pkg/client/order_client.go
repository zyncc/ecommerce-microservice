package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	pb "github.com/zyncc/ecommerce-microservice/services/order/pkg/types/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderClient interface {
	CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (uuid.UUID, error)
	FindOrderByOrderID(ctx context.Context, orderID uuid.UUID) (dto.FindOrderByIDResponse, error)
}

type OrderGRPCClient struct {
	client      pb.OrderServiceClient
	log         *zap.Logger
	orderSvcURL string
}

func NewOrderGRPCClient(log *zap.Logger, addr string) (*OrderGRPCClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &OrderGRPCClient{
		log:         log,
		client:      pb.NewOrderServiceClient(conn),
		orderSvcURL: addr,
	}, nil
}

func (c *OrderGRPCClient) CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (uuid.UUID, error) {
	var items []*pb.OrderItem
	for _, item := range req.Items {
		items = append(items, &pb.OrderItem{
			ProductId: item.ProductID.String(),
			Quantity:  int32(item.Quantity),
			Size:      item.Size,
			Price:     item.Price,
		})
	}

	resp, err := c.client.CreateOrder(ctx, &pb.CreateOrderRequest{
		Items:     items,
		UserId:    req.UserID.String(),
		AddressId: req.AddressID.String(),
	})
	if err != nil {
		c.log.Error("order service returned error", zap.Error(err))
		return uuid.Nil, err
	}

	orderID, err := uuid.Parse(resp.GetId())
	if err != nil {
		return uuid.Nil, err
	}

	return orderID, nil
}

func (c *OrderGRPCClient) FindOrderByOrderID(ctx context.Context, orderID uuid.UUID) (dto.FindOrderByIDResponse, error) {
	resp, err := c.client.FindOrderByID(ctx, &pb.IDMessage{Id: orderID.String()})
	if err != nil {
		c.log.Error("order service returned error", zap.Error(err))
		return dto.FindOrderByIDResponse{}, err
	}

	id, err := uuid.Parse(resp.GetId())
	if err != nil {
		return dto.FindOrderByIDResponse{}, err
	}

	userID, err := uuid.Parse(resp.GetUserId())
	if err != nil {
		return dto.FindOrderByIDResponse{}, err
	}

	orderItems := make([]dto.OrderItems, 0, len(resp.GetOrderItems()))
	for _, item := range resp.GetOrderItems() {
		itemID, err := uuid.Parse(item.GetId())
		if err != nil {
			return dto.FindOrderByIDResponse{}, err
		}

		orderItemOrderID, err := uuid.Parse(item.GetOrderId())
		if err != nil {
			return dto.FindOrderByIDResponse{}, err
		}

		productID, err := uuid.Parse(item.GetProductId())
		if err != nil {
			return dto.FindOrderByIDResponse{}, err
		}

		orderItems = append(orderItems, dto.OrderItems{
			ID:        itemID,
			OrderID:   orderItemOrderID,
			ProductID: productID,
			Quantity:  int(item.GetQuantity()),
			Size:      item.GetSize(),
			Price:     item.GetPrice(),
			CreatedAt: item.GetCreatedAt().AsTime(),
			UpdatedAt: item.GetUpdatedAt().AsTime(),
		})
	}

	return dto.FindOrderByIDResponse{
		ID:          id,
		UserID:      userID,
		Subtotal:    resp.GetSubtotal(),
		OrderTotal:  resp.GetOrderTotal(),
		OrderStatus: resp.GetOrderStatus(),
		FirstName:   resp.GetFirstName(),
		LastName:    resp.LastName,
		Email:       resp.GetEmail(),
		Phone:       resp.GetPhone(),
		Address1:    resp.GetAddress_1(),
		Address2:    resp.Address_2,
		City:        resp.GetCity(),
		State:       resp.GetState(),
		Zip:         resp.GetZip(),
		CreatedAt:   resp.GetCreatedAt().AsTime(),
		UpdatedAt:   resp.GetUpdatedAt().AsTime(),
		OrderItems:  orderItems,
	}, nil
}
