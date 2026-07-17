package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zyncc/ecommerce-microservice/services/payment/internal/repository/models"
	"github.com/zyncc/ecommerce-microservice/services/payment/pkg/types/topics"
	"go.uber.org/zap"
)

func (h *PaymentEventHandler) paymentSucceededEvent(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var event topics.PaymentSucceededEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		h.log.Error("invalid payment event, dropping", zap.Error(err), zap.ByteString("payload", msg.Value))
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := h.paymentRepo.CreatePayment(ctx, &models.CreatePaymentParams{
		OrderID:        event.OrderID,
		Status:         event.Status,
		Amount:         event.Amount,
		PaymentMethod:  event.PaymentMethod,
		Currency:       event.Currency,
		IdempotencyKey: event.IdempotencyKey,
	})
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" && pgErr.ConstraintName == "payments_order_id_key" {
			return nil
		}
		return err
	}

	return nil
}
