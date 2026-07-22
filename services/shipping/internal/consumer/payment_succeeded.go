package consumer

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/payment/pkg/types/topics"
	"github.com/zyncc/ecommerce-microservice/services/shipping/internal/repository/models"
	"go.uber.org/zap"
)

func (h *ShipmentEventHandler) paymentSucceededEvent(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var event topics.PaymentSucceededEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		h.log.Error("invalid payment event, dropping", zap.Error(err), zap.ByteString("payload", msg.Value))
		return nil
	}

	carriers := []string{"Blue Dart", "DTDC", "XpressBees", "Amazon Shipping", "FedEx", "DHL"}

	params := models.CreateShipmentParams{
		ID:             uuid.New(),
		Carrier:        carriers[rand.IntN(6)],
		OrderID:        event.OrderID,
		ShippingCost:   rand.Float64() * 100,
		TrackingNumber: uuid.New(),
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	_, err := h.shipmentRepo.CreateShipment(ctx, &params)
	if err != nil {
		cancel()
		return err
	}
	defer cancel()

	return nil
}
