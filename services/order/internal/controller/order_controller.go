package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/order/internal/service"
	pb "github.com/zyncc/ecommerce-microservice/services/order/pkg/types/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderController struct {
	pb.UnimplementedOrderServiceServer
	log *zap.Logger
	svc *service.OrderService
}

func NewOrderController(log *zap.Logger, svc *service.OrderService) *OrderController {
	return &OrderController{
		log: log,
		svc: svc,
	}
}

func (c *OrderController) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.IDMessage, error) {
	userID, err := uuid.Parse(in.GetUserId())
	if err != nil {
		c.log.Error("failed to parse user_id", zap.String("user_id", in.GetUserId()), zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	addressID, err := uuid.Parse(in.GetAddressId())
	if err != nil {
		c.log.Error("failed to parse address_id", zap.String("address_id", in.GetAddressId()), zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid address_id")
	}

	var items []dto.OrderItem
	for _, item := range in.GetItems() {
		productID, err := uuid.Parse(item.GetProductId())
		if err != nil {
			c.log.Error("failed to parse product_id", zap.String("product_id", item.GetProductId()), zap.Error(err))
			return nil, status.Error(codes.InvalidArgument, "invalid product_id")
		}

		items = append(items, dto.OrderItem{
			ProductID: productID,
			Quantity:  int(item.GetQuantity()),
			Size:      item.GetSize(),
			Price:     item.GetPrice(),
		})
	}

	req := dto.CreateOrderRequest{
		Items:     items,
		UserID:    userID,
		AddressID: addressID,
	}

	id, err := c.svc.CreateOrder(ctx, req)
	if err != nil {
		c.log.Error("failed to create order", zap.Any("request", req), zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create order")
	}

	return &pb.IDMessage{
		Id: id.String(),
	}, nil
}

func (c *OrderController) FindOrderByID(ctx context.Context, in *pb.IDMessage) (*pb.Order, error) {
	orderID, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "order id must be a valid uuid")
	}

	order, err := c.svc.FindOrderByOrderID(ctx, orderID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to find order by id")
	}

	var orderItems []*pb.OrderItems
	for _, item := range order.OrderItems {
		orderItems = append(orderItems, &pb.OrderItems{
			Id:        item.ID.String(),
			OrderId:   item.OrderID.String(),
			ProductId: item.ProductID.String(),
			Quantity:  int32(item.Quantity),
			Size:      item.Size,
			Price:     item.Price,
			CreatedAt: timestamppb.New(item.CreatedAt),
			UpdatedAt: timestamppb.New(item.UpdatedAt),
		})
	}

	return &pb.Order{
		Id:          order.ID.String(),
		UserId:      order.UserID.String(),
		Subtotal:    order.Subtotal,
		OrderTotal:  order.OrderTotal,
		OrderStatus: order.OrderStatus,
		FirstName:   order.FirstName,
		LastName:    order.LastName,
		Email:       order.Email,
		Phone:       order.Phone,
		Address_1:   order.Address1,
		Address_2:   order.Address2,
		City:        order.City,
		State:       order.State,
		Zip:         order.Zip,
		CreatedAt:   timestamppb.New(order.CreatedAt),
		UpdatedAt:   timestamppb.New(order.UpdatedAt),
		OrderItems:  orderItems,
	}, nil
}
