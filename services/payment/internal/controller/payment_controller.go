package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/payment/internal/service"
	"github.com/zyncc/ecommerce-microservice/services/payment/pkg/types/dto"
	pb "github.com/zyncc/ecommerce-microservice/services/payment/pkg/types/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PaymentController struct {
	pb.UnimplementedPaymentServiceServer
	log *zap.Logger
	svc *service.PaymentService
}

func NewPaymentController(log *zap.Logger, svc *service.PaymentService) *PaymentController {
	return &PaymentController{
		log: log,
		svc: svc,
	}
}

func (c *PaymentController) PaymentWebhook(ctx context.Context, in *pb.WebhookRequest) (*emptypb.Empty, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "razorpay header is missing")
	}

	// validate payment webhook signature to ensure it was sent from the payment provider
	razorpaySignature := meta.Get("X-Razorpay-Signature")[0]
	if razorpaySignature == "" {
		return nil, status.Error(codes.PermissionDenied, "you cannot access this resource")
	}

	idempotencyKey, _ := uuid.Parse(in.GetIdempotencyKey())
	orderID, _ := uuid.Parse(in.GetOrderId())

	if err := c.svc.ProcessPaymentWebhook(ctx, dto.PaymentWebhookRequest{
		IdempotencyKey: idempotencyKey,
		OrderID:        orderID,
		Amount:         in.GetAmount(),
		PaymentMethod:  in.GetPaymentMethod(),
		Currency:       in.GetCurrency(),
		Status:         in.GetStatus(),
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
