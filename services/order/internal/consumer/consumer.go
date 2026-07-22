package consumer

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/zyncc/ecommerce-microservice/services/order/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/order/pkg/types"
	"go.uber.org/zap"
)

type OrderEventHandler struct {
	orderRepo *repository.OrderRepository
	log       *zap.Logger
}

func NewOrderEventHandler(log *zap.Logger, orderRepo *repository.OrderRepository) *OrderEventHandler {
	return &OrderEventHandler{
		orderRepo,
		log,
	}
}

func (h *OrderEventHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *OrderEventHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *OrderEventHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			if err := h.processMessage(session.Context(), msg); err != nil {
				h.log.Error("failed to process payment event, will not commit offset", zap.Error(err), zap.Int32("partition", msg.Partition), zap.Int64("offset", msg.Offset))
				continue
			}
			session.MarkMessage(msg, "")
			h.log.Debug("marked offset for commit", zap.Int32("partition", msg.Partition), zap.Int64("offset", msg.Offset))

		case <-session.Context().Done():
			return nil
		}
	}
}

func (h *OrderEventHandler) processMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	switch msg.Topic {
	case types.PaymentSucceededTopic:
		if err := h.paymentSucceededEvent(ctx, msg); err != nil {
			return err
		}
	case types.ShipmentUpdatedTopic:
		if err := h.shipmentUpdatedEvent(ctx, msg); err != nil {
			return err
		}
	default:
		h.log.Error("unknown topic", zap.String("topic", msg.Topic))
	}

	return nil
}

type OrderConsumer struct {
	Brokers   []string
	Topics    []string
	GroupID   string
	OrderRepo *repository.OrderRepository
	Log       *zap.Logger
}

func (c *OrderConsumer) RunOrderConsumer(ctx context.Context) error {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Consumer.Offsets.AutoCommit.Enable = true
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup(c.Brokers, c.GroupID, saramaCfg)
	if err != nil {
		return err
	}
	defer group.Close()

	handler := NewOrderEventHandler(c.Log, c.OrderRepo)

	go func() {
		for err := range group.Errors() {
			c.Log.Error("kafka consumer group error", zap.Error(err))
		}
	}()

	for {
		c.Log.Info("Running Kafka Consumer")
		if err := group.Consume(ctx, c.Topics, handler); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}
			c.Log.Error("consumer group session error", zap.Error(err))
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}
