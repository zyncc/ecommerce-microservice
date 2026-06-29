package controller

import (
	"encoding/json"
	"net/http"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/product/internal/service"
	"github.com/zyncc/ecommerce-microservice/services/product/pkg/types/dto"
	"go.uber.org/zap"
)

type ProductController struct {
	log *zap.Logger
	svc *service.ProductService
}

func NewProductController(log *zap.Logger, svc *service.ProductService) *ProductController {
	return &ProductController{log, svc}
}

func (c *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	productId, err := c.svc.CreateProduct(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Product Created", &productId)
}

func (c *ProductController) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := c.svc.GetAllProducts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Fetched all products", &products)
}
