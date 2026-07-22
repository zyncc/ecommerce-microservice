package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/shipping/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/shipping/internal/repository/models"
	"github.com/zyncc/ecommerce-microservice/services/shipping/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/shipping/pkg/types/topics"
	"go.uber.org/zap"
)

type ShipmentService struct {
	log           *zap.Logger
	shipmentRepo  *repository.ShipmentRepository
	kafkaProducer sarama.SyncProducer
}

func NewShipmentService(log *zap.Logger, shipmentRepo *repository.ShipmentRepository, kafkaProducer sarama.SyncProducer) *ShipmentService {
	return &ShipmentService{log, shipmentRepo, kafkaProducer}
}

func (s *ShipmentService) ShipmentUpdateWebhook(ctx context.Context, req dto.ShipmentWebhookRequest) error {
	exists, err := s.shipmentRepo.FindShipmentByIdempotencyKey(ctx, req.IdempotencyKey)
	if err != nil {
		return err
	}

	if exists {
		s.log.Debug("webhook already processed")
		return nil
	}

	var shippedAt *time.Time
	var deliveredAt *time.Time

	now := time.Now()

	switch req.Status {
	case "SHIPPED":
		shippedAt = &now
	case "DELIVERED":
		deliveredAt = &now
	}

	params := models.UpdateShipmentParams{
		IdempotencyKey: req.IdempotencyKey,
		TrackingID:     req.TrackingNumber,
		Status:         req.Status,
		ShippedAt:      shippedAt,
		DeliveredAt:    deliveredAt,
	}

	if err := s.shipmentRepo.UpdateShipment(ctx, params); err != nil {
		return err
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: string(topics.ShipmentUpdatedTopic),
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := s.kafkaProducer.SendMessage(msg)
	if err != nil {
		return err
	}

	s.log.Info("Published event to kafka", zap.Int32("partition", partition), zap.Int64("offset", offset))

	return nil
}

func (s *ShipmentService) GetShipmentByTrackingID(ctx context.Context, trackingID uuid.UUID) (models.Shipment, error) {
	return s.shipmentRepo.GetShipmentByTrackingID(ctx, trackingID)
}
