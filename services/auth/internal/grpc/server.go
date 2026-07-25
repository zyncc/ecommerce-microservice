package grpc

import (
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/controller"
	pb "github.com/zyncc/ecommerce-microservice/services/auth/pkg/types/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(log *zap.Logger, authController *controller.AuthController) *grpc.Server {
	server := grpc.NewServer()

	pb.RegisterAuthServiceServer(server, authController)
	log.Info("registered gRPC services")

	return server
}
