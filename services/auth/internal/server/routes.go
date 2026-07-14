package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/controller"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/service"
	"github.com/zyncc/ecommerce-microservice/services/auth/pkg/middleware"
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
	userRepo := repository.NewUserRepository(s.log, s.pool)
	addressRepo := repository.NewAddressRepository(s.log, s.pool)

	// services
	authService := service.NewAuthService(s.log, userRepo, s.kafkaProducer, s.env)
	addressService := service.NewAddressService(s.log, addressRepo, s.env)

	// controllers
	authController := controller.NewAuthController(s.log, authService)
	addressController := controller.NewAddressController(s.log, addressService)

	// routes
	r.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		utils.SuccessResponse[any](w, http.StatusOK, "auth service healthy", nil)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.HandleFunc("POST /signup", authController.SignUp)
		r.HandleFunc("POST /signin", authController.SignIn)
		r.HandleFunc("POST /refresh", authController.RefreshToken)
		r.HandleFunc("GET /session", authController.GetSession)

		r.HandleFunc("POST /address", addressController.CreateAddress)
		r.HandleFunc("GET /address/{id}", addressController.GetAddressByID)
		r.HandleFunc("GET /address", addressController.FetchAllAddresses)
	})

	return r
}
