package controller

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/service"
	"go.uber.org/zap"
)

type AddressController struct {
	logger *zap.Logger
	svc    *service.AddressService
}

func NewAddressController(logger *zap.Logger, svc *service.AddressService) *AddressController {
	return &AddressController{
		logger,
		svc,
	}
}

func (c *AddressController) CreateAddress(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.logger.Debug("failed to parse request body", zap.Error(err))
		utils.ErrorResponse(w, http.StatusBadRequest, "failed to parse request body")
		return
	}

	id, err := c.svc.CreateAddress(r.Context(), req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "failed to create address")
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Created Address", id)
}

func (c *AddressController) GetAddressByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "address id is required")
		return
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "id is not a valid uuid")
		return
	}

	address, err := c.svc.FindAddressByID(r.Context(), parsedID)
	if err != nil {
		utils.ErrorResponse(w, 500, "failed to find address")
		return
	}

	utils.SuccessResponse(w, 200, "fetched address by id", address)
}

func (c *AddressController) FetchAllAddresses(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")

	if userID == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "user id is required")
		return
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "uuid is not a valid uuid")
		return
	}

	addresses, err := c.svc.FetchAllAddresses(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "failed to fetch addresses")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Fetched all addresses", addresses)
}
