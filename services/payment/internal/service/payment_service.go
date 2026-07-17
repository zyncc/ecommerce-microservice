package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/payment/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/payment/pkg/types"
	"github.com/zyncc/ecommerce-microservice/services/payment/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/payment/pkg/types/topics"
	"go.uber.org/zap"
)

type PaymentService struct {
	log           *zap.Logger
	orderRepo     *repository.OrderRepository
	kafkaProducer sarama.SyncProducer
}

func NewPaymentService(log *zap.Logger, orderRepo *repository.OrderRepository, kafkaProducer sarama.SyncProducer) *PaymentService {
	return &PaymentService{log, orderRepo, kafkaProducer}
}

func (s *PaymentService) ProcessPaymentWebhook(ctx context.Context, req dto.PaymentWebhookRequest) error {
	// check for idempotency key to prevent duplicate processing
	exists, err := s.orderRepo.FindOrderByIdempotencyKey(ctx, req.IdempotencyKey)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	event := topics.PaymentSucceededEvent{
		EventID:        uuid.New(),
		IdempotencyKey: req.IdempotencyKey,
		Amount:         req.Amount,
		OrderID:        req.OrderID,
		Status:         req.Status,
		OccurredAt:     time.Now(),
		PaymentMethod:  req.PaymentMethod,
		Currency:       req.Currency,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(req.OrderID.String()),
		Topic: types.PaymentSucceededTopic,
		Value: sarama.ByteEncoder(payload),
	}

	partition, offset, err := s.kafkaProducer.SendMessage(msg)
	if err != nil {
		s.log.Error("failed to send kafka message", zap.Error(err))
		return err
	}

	s.log.Info("published payment success kafka message", zap.Int32("partition", partition), zap.Int64("offset", offset))

	return nil
}
