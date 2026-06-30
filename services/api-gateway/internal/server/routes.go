package server

import (
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/zyncc/ecommerce-microservice/services/api-gateway/docs"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/internal/client"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/internal/controller"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/middleware"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger(s.log))
	r.Use(chimiddleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	// clients
	authClient := client.NewAuthClient(s.log, s.env, httpClient)
	productClient := client.NewProductClient(s.log, s.env, httpClient)

	// controller
	authController := controller.NewAuthController(s.log, authClient)
	productController := controller.NewProductController(s.log, productClient)

	// middleware
	authMiddleware := middleware.NewAuthMiddleware(s.log, authClient)

	r.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		utils.SuccessResponse[any](w, http.StatusOK, "api gateway healthy", nil)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	r.Route("/api/v1", func(r chi.Router) {
		r.HandleFunc("POST /signup", authController.SignUp)
		r.HandleFunc("POST /signin", authController.SignIn)
		r.HandleFunc("POST /refresh", authController.RefreshToken)

		r.HandleFunc("GET /product", productController.GetAllProducts)
		r.HandleFunc("GET /product/{id}", productController.GetProductByID)

		// authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.HandleFunc("GET /session", authController.GetSession)
			r.HandleFunc("POST /signout", authController.SignOut)
		})

		// admin routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAdmin)
			r.HandleFunc("POST /product", productController.CreateProduct)
		})
	})

	return r
}
