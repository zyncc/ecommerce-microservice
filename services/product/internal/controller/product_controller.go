package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/product/internal/service"
	"github.com/zyncc/ecommerce-microservice/services/product/pkg/types"
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

	productID, err := c.svc.CreateProduct(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Product Created", productID)
}

const (
	defaultLimit = 5
	maxLimit     = 10
)

func (s *ProductController) GetAllProducts(w http.ResponseWriter, r *http.Request) {
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

	offset := 0
	if offsetQuery := r.URL.Query().Get("offset"); offsetQuery != "" {
		parsed, err := strconv.Atoi(offsetQuery)
		if err != nil || parsed < 0 {
			utils.ErrorResponse(w, http.StatusBadRequest, "offset needs to be a positive number")
			return
		}
		offset = parsed
	}

	products, err := s.svc.GetAllProducts(r.Context(), limit, offset)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(w, http.StatusOK, "Fetched all products", &products)
}

func (c *ProductController) GetProductByID(w http.ResponseWriter, r *http.Request) {
	pathID := r.PathValue("id")

	id, err := uuid.Parse(pathID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "id is not valid")
		return
	}

	productID, err := c.svc.GetProductByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, types.ErrProductNotFound) {
			utils.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Fetched Product", &productID)
}

func (c *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	pathID := r.PathValue("id")

	id, err := uuid.Parse(pathID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "id is not valid")
		return
	}

	if err := c.svc.DeleteProduct(r.Context(), id); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse[any](w, http.StatusOK, "Deleted Product", nil)
}
