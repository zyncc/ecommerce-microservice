package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/product/pkg/middleware"
	"github.com/zyncc/ecommerce-microservice/services/shipping/internal/controller"
	"github.com/zyncc/ecommerce-microservice/services/shipping/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/shipping/internal/service"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger(s.log))
	r.Use(chimiddleware.Recoverer)

	// repositories
	shipmentRepo := repository.NewShippingRepository(s.log, s.pool)

	// services
	shipmentService := service.NewShipmentService(s.log, shipmentRepo, s.kafkaProducer)

	// controllers
	shipmentController := controller.NewShipmentController(s.log, shipmentService)

	// routes
	r.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		utils.SuccessResponse[any](w, http.StatusOK, "shipping service healthy", nil)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.HandleFunc("POST /webhook/shipment", shipmentController.ShipmentUpdateWebhook)
		r.HandleFunc("GET /shipment", shipmentController.GetShipmentByTrackingID)
	})

	return r
}
