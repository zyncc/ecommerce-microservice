package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
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
		if httpErr, ok := errors.AsType[*utils.HTTPError](err); ok {
			utils.ErrorResponse(w, httpErr.Status, httpErr.Message)
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
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
	const (
		defaultLimit = 5
		maxLimit     = 100
		defaultPage  = 1
	)

	limit := defaultLimit
	if limitQuery := r.URL.Query().Get("limit"); limitQuery != "" {
		parsed, err := strconv.Atoi(limitQuery)
		if err != nil || parsed <= 0 {
			utils.ErrorResponse(w, http.StatusBadRequest, "limit needs to be a positive number")
			return
		}
		if parsed > maxLimit {
			parsed = maxLimit
		}
		limit = parsed
	}

	page := defaultPage
	if pageQuery := r.URL.Query().Get("page"); pageQuery != "" {
		parsed, err := strconv.Atoi(pageQuery)
		if err != nil || parsed <= 0 {
			utils.ErrorResponse(w, http.StatusBadRequest, "page needs to be a positive number")
			return
		}
		page = parsed
	}

	offset := (page - 1) * limit

	products, err := c.productClient.GetAllProducts(r.Context(), limit, offset)
	if err != nil {
		if httpErr, ok := errors.AsType[*utils.HTTPError](err); ok {
			utils.ErrorResponse(w, httpErr.Status, httpErr.Message)
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	utils.SuccessResponse(w, http.StatusOK, "Products Retrieved", &products)
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
func (c *ProductController) GetProductByID(w http.ResponseWriter, r *http.Request) {
	pathID := r.PathValue("id")
	id, err := uuid.Parse(pathID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "id is not valid")
		return
	}

	resp, err := c.productClient.GetProductByID(r.Context(), id)
	if err != nil {
		if httpErr, ok := errors.AsType[*utils.HTTPError](err); ok {
			utils.ErrorResponse(w, httpErr.Status, httpErr.Message)
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Fetched Product", resp.Data)
}
