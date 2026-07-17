package consumer

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/zyncc/ecommerce-microservice/services/payment/pkg/types/topics"
	"go.uber.org/zap"
)

func (h *InventoryEventHandler) paymentSucceededEvent(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var event topics.PaymentSucceededEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		h.log.Error("invalid payment event, dropping", zap.Error(err), zap.ByteString("payload", msg.Value))
		return nil
	}

	return nil
}
