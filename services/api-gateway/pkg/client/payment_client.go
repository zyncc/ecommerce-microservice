package client

import (
	"context"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	pb "github.com/zyncc/ecommerce-microservice/services/payment/pkg/types/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type PaymentClient interface {
	PaymentWebhook(ctx context.Context, req dto.PaymentWebhookRequest, razorpaySignature string) error
}

type PaymentGRPCClient struct {
	log           *zap.Logger
	client        pb.PaymentServiceClient
	paymentSvcURL string
}

func NewPaymentGRPCClient(log *zap.Logger, addr string) (*PaymentGRPCClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &PaymentGRPCClient{
		log:           log,
		paymentSvcURL: addr,
		client:        pb.NewPaymentServiceClient(conn),
	}, nil
}

func (c *PaymentGRPCClient) PaymentWebhook(ctx context.Context, req dto.PaymentWebhookRequest, razorpaySignature string) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "X-Razorpay-Signature", razorpaySignature)
	_, err := c.client.PaymentWebhook(ctx, &pb.WebhookRequest{
		OrderId:        req.OrderID.String(),
		IdempotencyKey: req.IdempotencyKey.String(),
		Amount:         req.Amount,
		PaymentMethod:  req.PaymentMethod,
		Status:         req.Status,
		Currency:       req.Currency,
	})
	if err != nil {
		return err
	}

	return nil
}
