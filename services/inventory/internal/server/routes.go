package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/controller"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/service"
	"github.com/zyncc/ecommerce-microservice/services/product/pkg/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger(s.log))
	r.Use(chimiddleware.Recoverer)

	// repository
	inventoryRepo := repository.NewInventoryRepository(s.log, s.pool)

	// services
	inventoryService := service.NewProductService(s.log, inventoryRepo)

	// controllers
	inventoryController := controller.NewInventoryController(s.log, inventoryService)

	// routes
	r.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		utils.SuccessResponse[any](w, http.StatusOK, "inventory service healthy", nil)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.HandleFunc("POST /inventory", inventoryController.CreateInventory)
		r.HandleFunc("PUT /inventory", inventoryController.UpdateInventory)
		r.HandleFunc("GET /inventory/{productID}", inventoryController.GetInventoryByProductID)
	})

	return r
}
