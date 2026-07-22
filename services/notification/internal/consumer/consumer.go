package consumer

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/client"
	"github.com/zyncc/ecommerce-microservice/services/notification/internal/config"
	"github.com/zyncc/ecommerce-microservice/services/notification/pkg/types"
	"go.uber.org/zap"
)

type NotificationEventHandler struct {
	log         *zap.Logger
	orderClient *client.OrderClient
	env         *config.EnvConfig
}

func NewNotificationEventHandler(log *zap.Logger, orderClient *client.OrderClient, env *config.EnvConfig) *NotificationEventHandler {
	return &NotificationEventHandler{
		log,
		orderClient,
		env,
	}
}

func (h *NotificationEventHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *NotificationEventHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *NotificationEventHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
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

func (h *NotificationEventHandler) processMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
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

type NotificationConsumer struct {
	Brokers     []string
	OrderClient *client.OrderClient
	Env         *config.EnvConfig
	Topics      []string
	GroupID     string
	Log         *zap.Logger
}

func (c *NotificationConsumer) RunNotificationConsumer(ctx context.Context) error {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Consumer.Offsets.AutoCommit.Enable = true
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup(c.Brokers, c.GroupID, saramaCfg)
	if err != nil {
		return err
	}
	defer group.Close()

	handler := NewNotificationEventHandler(c.Log, c.OrderClient, c.Env)

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
