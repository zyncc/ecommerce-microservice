package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/payment/internal/controller"
	"github.com/zyncc/ecommerce-microservice/services/payment/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/payment/internal/service"
	"github.com/zyncc/ecommerce-microservice/services/product/pkg/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger(s.log))
	r.Use(chimiddleware.Recoverer)

	// repositories
	paymentRepo := repository.NewPaymentRepository(s.log, s.pool)

	// services
	paymentService := service.NewPaymentService(s.log, paymentRepo, s.kafkaProducer)

	// controllers
	paymentController := controller.NewPaymentController(s.log, paymentService)

	// routes
	r.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		utils.SuccessResponse[any](w, http.StatusOK, "inventory service healthy", nil)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.HandleFunc("POST /webhook/payment", paymentController.PaymentWebhook)
	})

	return r
}
