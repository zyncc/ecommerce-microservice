package consumer

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/zyncc/ecommerce-microservice/services/payment/pkg/types"
	"github.com/zyncc/ecommerce-microservice/services/shipping/internal/repository"
	"go.uber.org/zap"
)

type ShipmentEventHandler struct {
	shipmentRepo *repository.ShipmentRepository
	log          *zap.Logger
}

func NewShipmentEventHandler(log *zap.Logger, shipmentRepo *repository.ShipmentRepository) *ShipmentEventHandler {
	return &ShipmentEventHandler{
		shipmentRepo,
		log,
	}
}

func (h *ShipmentEventHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ShipmentEventHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ShipmentEventHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
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

func (h *ShipmentEventHandler) processMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	switch msg.Topic {
	case types.PaymentSucceededTopic:
		if err := h.paymentSucceededEvent(ctx, msg); err != nil {
			return err
		}
	default:
		h.log.Error("unknown topic", zap.String("topic", msg.Topic))
	}

	return nil
}

type ShipmentConsumer struct {
	Brokers      []string
	Topics       []string
	GroupID      string
	ShipmentRepo *repository.ShipmentRepository
	Log          *zap.Logger
}

func (c *ShipmentConsumer) RunShipmentConsumer(ctx context.Context) error {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Consumer.Offsets.AutoCommit.Enable = true
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup(c.Brokers, c.GroupID, saramaCfg)
	if err != nil {
		return err
	}
	defer group.Close()

	handler := NewShipmentEventHandler(c.Log, c.ShipmentRepo)

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
