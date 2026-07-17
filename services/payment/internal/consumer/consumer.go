package consumer

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/zyncc/ecommerce-microservice/services/payment/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/payment/pkg/types"
	"go.uber.org/zap"
)

type PaymentEventHandler struct {
	paymentRepo *repository.PaymentRepository
	log         *zap.Logger
}

func NewPaymentEventHandler(log *zap.Logger, paymentRepo *repository.PaymentRepository) *PaymentEventHandler {
	return &PaymentEventHandler{
		paymentRepo: paymentRepo,
		log:         log,
	}
}

func (h *PaymentEventHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *PaymentEventHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *PaymentEventHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
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

func (h *PaymentEventHandler) processMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
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

type PaymentConsumer struct {
	Brokers     []string
	Topics      []string
	GroupID     string
	PaymentRepo *repository.PaymentRepository
	Log         *zap.Logger
}

func (c *PaymentConsumer) RunPaymentConsumer(ctx context.Context) error {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Consumer.Offsets.AutoCommit.Enable = true
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup(c.Brokers, c.GroupID, saramaCfg)
	if err != nil {
		return err
	}
	defer group.Close()

	handler := NewPaymentEventHandler(c.Log, c.PaymentRepo)

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
