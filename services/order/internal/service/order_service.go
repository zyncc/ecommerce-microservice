package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/client"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/order/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/order/internal/repository/model"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type OrderService struct {
	log             *zap.Logger
	repo            *repository.OrderRepository
	authClient      *client.AuthClient
	productClient   *client.ProductClient
	inventoryClient *client.InventoryClient
}

func NewOrderService(log *zap.Logger, repo *repository.OrderRepository, authClient *client.AuthClient, productClient *client.ProductClient, inventoryClient *client.InventoryClient) *OrderService {
	return &OrderService{log, repo, authClient, productClient, inventoryClient}
}

func (s *OrderService) CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (uuid.UUID, error) {
	// fetch address
	// compute order price
	// check if inventory exists

	address, err := s.authClient.GetAddressByID(ctx, req.AddressID)
	if err != nil {
		return uuid.Nil, err
	}

	// compute order price
	var (
		subtotal float64
		mu       sync.Mutex
	)

	g, groupCtx := errgroup.WithContext(ctx)
	g.SetLimit(10)
	for i := range req.Items {
		g.Go(func() error {
			product, err := s.productClient.GetProductByID(groupCtx, req.Items[i].ProductID)
			if err != nil {
				return err
			}

			req.Items[i].Price = product.Price
			mu.Lock()
			subtotal += product.Price * float64(req.Items[i].Quantity)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return uuid.Nil, err
	}

	// fetch shipping cost
	// shippingCost, err := s.shippingClient.GetShippingCost(ctx, address.Zip)
	// if err != nil {
	// 	return uuid.Nil, err
	// }

	shippingCost := 50 + rand.Float64()*150
	shippingCost = math.Round(shippingCost*100) / 100
	orderTotal := subtotal + shippingCost

	// check inventory
	g, groupCtx = errgroup.WithContext(ctx)
	g.SetLimit(10)
	for i := range req.Items {
		g.Go(func() error {
			inventory, err := s.inventoryClient.FetchInventoryByProductID(groupCtx, req.Items[i].ProductID)
			if err != nil {
				return err
			}

			var available int

			switch req.Items[i].Size {
			case "small":
				available = inventory.Small
			case "medium":
				available = inventory.Medium
			case "large":
				available = inventory.Large
			case "extra_large":
				available = inventory.ExtraLarge
			default:
				return errors.New("invalid size")
			}

			if available < req.Items[i].Quantity {
				return fmt.Errorf("inventory out of stock for size %s", req.Items[i].Size)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return uuid.Nil, err
	}

	id, err := s.repo.CreateOrder(ctx, model.CreateOrderParams{
		UserID:       req.UserID,
		Items:        req.Items,
		Subtotal:     subtotal,
		OrderTotal:   orderTotal,
		ShippingCost: shippingCost,
		FirstName:    address.FirstName,
		LastName:     address.LastName,
		Email:        address.Email,
		Phone:        address.Phone,
		Address1:     address.Address1,
		Address2:     address.Address2,
		City:         address.City,
		State:        address.State,
		Zip:          address.Zip,
	})
	if err != nil {
		s.log.Error("failed to create order", zap.Error(err))
		return uuid.Nil, err
	}

	return id, nil
}
