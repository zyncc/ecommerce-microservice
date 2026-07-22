package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/zyncc/ecommerce-microservice/services/order/pkg/types"
	"go.uber.org/zap"
)

func (h *OrderEventHandler) shipmentUpdatedEvent(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var event types.ShipmentWebhookRequest
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		h.log.Error("invalid payment event, dropping", zap.Error(err), zap.ByteString("payload", msg.Value))
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	if err := h.orderRepo.UpdateOrderStatusWithTrackingID(ctx, event.TrackingNumber, event.Status); err != nil {
		cancel()
		return err
	}
	defer cancel()

	return nil
}
