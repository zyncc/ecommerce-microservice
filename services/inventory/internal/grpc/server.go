package grpc

import (
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/controller"
	pb "github.com/zyncc/ecommerce-microservice/services/inventory/pkg/types/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(log *zap.Logger, inventoryController *controller.InventoryController) *grpc.Server {
	server := grpc.NewServer()

	pb.RegisterInventoryServiceServer(server, inventoryController)
	log.Info("registered gRPC services")

	return server
}
