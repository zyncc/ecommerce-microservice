package controller

import (
	"encoding/json"
	"net/http"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/internal/client"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/product/pkg/types/dto"
	"go.uber.org/zap"
)

type ProductController struct {
	log           *zap.Logger
	productClient *client.ProductClient
}

func NewProductController(log *zap.Logger, productClient *client.ProductClient) *ProductController {
	return &ProductController{
		log,
		productClient,
	}
}

// CreateProduct godoc
// @Summary Create Product
// @Description Creates a new Product
// @Tags Product
// @Accept json
// @Produce json
// @Param request body dto.CreateProductRequest true "Create Product Request"
// @Success 200 {object} uuid.UUID
// @Failure 500 {object} utils.Error
// @Router /api/v1/product [post]
func (c *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid or malformed request")
		return
	}

	id, err := c.productClient.CreateProduct(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(w, http.StatusOK, "Product Created", &id)
}

// GetAllProducts godoc
// @Summary Fetch all products
// @Description Fetch all products
// @Tags Product
// @Accept json
// @Produce json
// @Success 200 {object} []dto.Product
// @Failure 500 {object} utils.Error
// @Router /api/v1/product [get]
func (c *ProductController) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := c.productClient.GetAllProducts(r.Context())
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(w, http.StatusOK, "Products Retrieved", &products)
}
