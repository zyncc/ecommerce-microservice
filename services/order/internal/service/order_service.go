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
	address, err := s.authClient.GetAddressByID(ctx, req.AddressID)
	if err != nil {
		return uuid.Nil, err
	}

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
		UserID:     req.UserID,
		Items:      req.Items,
		Subtotal:   subtotal,
		OrderTotal: orderTotal,
		FirstName:  address.FirstName,
		LastName:   address.LastName,
		Email:      address.Email,
		Phone:      address.Phone,
		Address1:   address.Address1,
		Address2:   address.Address2,
		City:       address.City,
		State:      address.State,
		Zip:        address.Zip,
	})
	if err != nil {
		s.log.Error("failed to create order", zap.Error(err))
		return uuid.Nil, err
	}

	return id, nil
}

func (s *OrderService) FindOrderByOrderID(ctx context.Context, orderID uuid.UUID) (dto.FindOrderByIDResponse, error) {
	order, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return dto.FindOrderByIDResponse{}, err
	}

	var orderItemsResp []dto.OrderItems

	for _, item := range order.OrderItems {
		respItem := dto.OrderItems{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Size:      item.Size,
			Price:     item.Price,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}

		orderItemsResp = append(orderItemsResp, respItem)
	}

	orderResp := dto.FindOrderByIDResponse{
		ID:             order.ID,
		UserID:         order.UserID,
		IdempotencyKey: order.IdempotencyKey,
		Subtotal:       order.Subtotal,
		OrderTotal:     order.OrderTotal,
		OrderStatus:    order.OrderStatus,
		FirstName:      order.FirstName,
		LastName:       order.LastName,
		Email:          order.Email,
		Phone:          order.Phone,
		Address1:       order.Address1,
		Address2:       order.Address2,
		City:           order.City,
		State:          order.State,
		Zip:            order.Zip,
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
		OrderItems:     orderItemsResp,
	}

	return orderResp, nil
}
