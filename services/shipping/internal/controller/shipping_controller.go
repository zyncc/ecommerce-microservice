package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/shipping/internal/service"
	"github.com/zyncc/ecommerce-microservice/services/shipping/pkg/types/dto"
	pb "github.com/zyncc/ecommerce-microservice/services/shipping/pkg/types/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ShipmentController struct {
	pb.UnimplementedShippingServiceServer
	log *zap.Logger
	svc *service.ShipmentService
}

func NewShipmentController(log *zap.Logger, svc *service.ShipmentService) *ShipmentController {
	return &ShipmentController{
		log: log,
		svc: svc,
	}
}

func (c *ShipmentController) ShipmentWebhook(ctx context.Context, in *pb.WebhookRequest) (*emptypb.Empty, error) {
	idempotencyKey, _ := uuid.Parse(in.GetIdempotencyKey())
	trackingNumber, _ := uuid.Parse(in.GetTrackingNumber())

	if err := c.svc.ShipmentUpdateWebhook(ctx, dto.ShipmentWebhookRequest{
		IdempotencyKey: idempotencyKey,
		TrackingNumber: trackingNumber,
		Status:         in.GetStatus(),
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
func (c *ShipmentController) GetShipmentByTrackingID(ctx context.Context, in *pb.IDMessage) (*pb.Shipment, error) {
	trackingID, _ := uuid.Parse(in.GetId())
	shipment, err := c.svc.GetShipmentByTrackingID(ctx, trackingID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	idempotencyKey := shipment.IdempotencyKey.String()
	shippedAt := timestamppb.New(*shipment.ShippedAt)
	deliveredAt := timestamppb.New(*shipment.DeliveredAt)
	return &pb.Shipment{
		Id:             shipment.ID.String(),
		OrderId:        shipment.OrderID.String(),
		IdempotencyKey: &idempotencyKey,
		Status:         shipment.Status,
		Carrier:        shipment.Carrier,
		ShippingCost:   shipment.ShippingCost,
		TrackingNumber: shipment.TrackingNumber.String(),
		ShippedAt:      shippedAt,
		DeliveredAt:    deliveredAt,
		CreatedAt:      timestamppb.New(shipment.CreatedAt),
		UpdatedAt:      timestamppb.New(shipment.UpdatedAt),
	}, nil
}
