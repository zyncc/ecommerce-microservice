package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/client"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/order/internal/controller"
	"github.com/zyncc/ecommerce-microservice/services/order/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/order/internal/service"
	"github.com/zyncc/ecommerce-microservice/services/product/pkg/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger(s.log))
	r.Use(chimiddleware.Recoverer)

	httpClient := http.Client{
		Timeout: time.Second * 5,
	}
	// clients
	authClient := client.NewAuthClient(s.log, s.env.AuthServiceURL, &httpClient)
	productClient := client.NewProductClient(s.log, s.env.ProductServiceURL, &httpClient)
	inventoryClient := client.NewInventoryClient(s.log, s.env.InventoryServiceURL, &httpClient)

	// repository
	orderRepo := repository.NewOrderRepository(s.log, s.pool)

	// services
	orderService := service.NewOrderService(s.log, orderRepo, authClient, productClient, inventoryClient)

	// controllers
	orderController := controller.NewOrderController(s.log, orderService)

	// routes
	r.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		utils.SuccessResponse[any](w, http.StatusOK, "order service healthy", nil)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.HandleFunc("POST /order", orderController.CreateOrder)
	})

	return r
}
