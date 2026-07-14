package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/config"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/repository/models"
	"go.uber.org/zap"
)

type AddressService struct {
	logger      *zap.Logger
	addressRepo *repository.AddressRepository
	env         *config.EnvConfig
}

func NewAddressService(logger *zap.Logger, addressRepo *repository.AddressRepository, env *config.EnvConfig) *AddressService {
	return &AddressService{
		logger,
		addressRepo,
		env,
	}
}

func (s *AddressService) CreateAddress(ctx context.Context, req dto.CreateAddressRequest) (uuid.UUID, error) {
	id := uuid.New()
	err := s.addressRepo.CreateAddress(ctx, &models.CreateAddressParams{
		ID:        id,
		UserID:    req.UserID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Address1:  req.Address1,
		Address2:  req.Address2,
		City:      req.City,
		State:     req.State,
		Zip:       req.Zip,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *AddressService) FindAddressByID(ctx context.Context, id uuid.UUID) (*models.Address, error) {
	address, err := s.addressRepo.FindAddressByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (s *AddressService) FetchAllAddresses(ctx context.Context, userID uuid.UUID) ([]models.Address, error) {
	addresses, err := s.addressRepo.FetchAllAddresses(ctx, userID)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}
