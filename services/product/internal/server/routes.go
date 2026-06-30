package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/product/internal/controller"
	"github.com/zyncc/ecommerce-microservice/services/product/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/product/internal/service"
	"github.com/zyncc/ecommerce-microservice/services/product/pkg/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger(s.log))
	r.Use(chimiddleware.Recoverer)

	// repository
	productRepo := repository.NewProductRepository(s.log, s.pool)

	// services
	productService := service.NewProductService(s.log, productRepo)

	// controllers
	productController := controller.NewProductController(s.log, productService)

	// routes
	r.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		utils.SuccessResponse[any](w, http.StatusOK, "product service healthy", nil)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.HandleFunc("POST /product", productController.CreateProduct)
		r.HandleFunc("GET /product", productController.GetAllProducts)
		r.HandleFunc("GET /product/{id}", productController.GetProductByID)
	})

	return r
}
