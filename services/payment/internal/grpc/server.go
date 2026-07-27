package grpc

import (
	"github.com/zyncc/ecommerce-microservice/services/payment/internal/controller"
	pb "github.com/zyncc/ecommerce-microservice/services/payment/pkg/types/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(log *zap.Logger, paymentController *controller.PaymentController) *grpc.Server {
	server := grpc.NewServer()

	pb.RegisterPaymentServiceServer(server, paymentController)
	log.Info("registered gRPC services")

	return server
}
