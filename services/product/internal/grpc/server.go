package grpc

import (
	"github.com/zyncc/ecommerce-microservice/services/product/internal/controller"
	pb "github.com/zyncc/ecommerce-microservice/services/product/pkg/types/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(log *zap.Logger, productController *controller.ProductController) *grpc.Server {
	server := grpc.NewServer()

	pb.RegisterProductServiceServer(server, productController)
	log.Info("registered gRPC services")

	return server
}
