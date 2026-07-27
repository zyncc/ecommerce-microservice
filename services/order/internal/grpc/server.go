package grpc

import (
	"github.com/zyncc/ecommerce-microservice/services/order/internal/controller"
	pb "github.com/zyncc/ecommerce-microservice/services/order/pkg/types/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(log *zap.Logger, orderController *controller.OrderController) *grpc.Server {
	server := grpc.NewServer()

	pb.RegisterOrderServiceServer(server, orderController)
	log.Info("registered gRPC services")

	return server
}
