package consumer

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/zyncc/ecommerce-microservice/services/inventory/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/payment/pkg/types/topics"
	"go.uber.org/zap"
)

func (h *InventoryEventHandler) paymentSucceededEvent(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var event topics.PaymentSucceededEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		h.log.Error("invalid payment event, dropping", zap.Error(err), zap.ByteString("payload", msg.Value))
		return nil
	}

	order, err := h.orderclient.FindOrderByOrderID(ctx, event.OrderID)
	if err != nil {
		return err
	}

	var inventoryParams []dto.UpdateInventoryRequest
	for _, item := range order.OrderItems {
		inventoryParam := dto.UpdateInventoryRequest{
			ProductID: item.ProductID,
			Size:      item.Size,
			Quantity:  item.Quantity,
		}

		inventoryParams = append(inventoryParams, inventoryParam)
	}

	if err := h.inventoryRepo.UpdateInventory(ctx, inventoryParams); err != nil {
		return err
	}

	return nil
}
