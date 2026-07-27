package grpc

import (
	"github.com/zyncc/ecommerce-microservice/services/shipping/internal/controller"
	pb "github.com/zyncc/ecommerce-microservice/services/shipping/pkg/types/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(log *zap.Logger, shipmentController *controller.ShipmentController) *grpc.Server {
	server := grpc.NewServer()

	pb.RegisterShippingServiceServer(server, shipmentController)
	log.Info("registered gRPC services")

	return server
}
