package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	pb "github.com/zyncc/ecommerce-microservice/services/shipping/pkg/types/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ShipmentClient interface {
	ShipmentWebhook(ctx context.Context, req dto.ShipmentWebhookRequest, shipmentSignature string) error
	GetShipmentByTrackingID(ctx context.Context, trackingID uuid.UUID) (dto.ShipmentResponse, error)
}

type ShipmentGRPCClient struct {
	log            *zap.Logger
	shipmentSvcURL string
	client         pb.ShippingServiceClient
}

func NewShipmentGRPCClient(log *zap.Logger, addr string) (*ShipmentGRPCClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &ShipmentGRPCClient{
		log:            log,
		shipmentSvcURL: addr,
		client:         pb.NewShippingServiceClient(conn),
	}, nil
}

func (c *ShipmentGRPCClient) ShipmentWebhook(ctx context.Context, req dto.ShipmentWebhookRequest, shipmentSignature string) error {
	_, err := c.client.ShipmentWebhook(ctx, &pb.WebhookRequest{
		IdempotencyKey: req.IdempotencyKey.String(),
		TrackingNumber: req.TrackingNumber.String(),
		Status:         req.Status,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *ShipmentGRPCClient) GetShipmentByTrackingID(ctx context.Context, trackingID uuid.UUID) (dto.ShipmentResponse, error) {
	resp, err := c.client.GetShipmentByTrackingID(ctx, &pb.IDMessage{
		Id: trackingID.String(),
	})
	if err != nil {
		return dto.ShipmentResponse{}, err
	}

	id, _ := uuid.Parse(resp.GetId())
	orderID, _ := uuid.Parse(resp.GetOrderId())
	idempotencyKey, _ := uuid.Parse(resp.GetIdempotencyKey())
	trackingNumber, _ := uuid.Parse(resp.GetTrackingNumber())

	shippedAt := resp.GetShippedAt().AsTime()
	deliveredAt := resp.GetDeliveredAt().AsTime()

	return dto.ShipmentResponse{
		ID:             id,
		OrderID:        orderID,
		IdempotencyKey: &idempotencyKey,
		Status:         resp.GetStatus(),
		Carrier:        resp.GetCarrier(),
		ShippingCost:   resp.GetShippingCost(),
		TrackingNumber: trackingNumber,
		ShippedAt:      &shippedAt,
		DeliveredAt:    &deliveredAt,
		CreatedAt:      resp.GetCreatedAt().AsTime(),
		UpdatedAt:      resp.GetUpdatedAt().AsTime(),
	}, nil
}
