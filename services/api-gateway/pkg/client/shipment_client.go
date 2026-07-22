package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"go.uber.org/zap"
)

type ShipmentClient struct {
	log            *zap.Logger
	shipmentSvcURL string
	httpClient     *http.Client
}

func NewShipmentClient(log *zap.Logger, shipmentSvcURL string, httpClient *http.Client) *ShipmentClient {
	return &ShipmentClient{
		log,
		shipmentSvcURL,
		httpClient,
	}
}

func (c *ShipmentClient) ShipmentWebhook(ctx context.Context, req dto.ShipmentWebhookRequest, shipmentSignature string) error {
	reqBody, err := json.Marshal(req)
	if err != nil {
		c.log.Error("failed to marshal json data", zap.Error(err))
		return utils.ErrSomethingWentWrong
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/v1/webhook/shipment", c.shipmentSvcURL), bytes.NewReader(reqBody))
	if err != nil {
		c.log.Error("failed to create http request", zap.Error(err))
		return utils.ErrSomethingWentWrong
	}

	request.Header.Set("X-Shipment-Signature", shipmentSignature)
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Error("failed to send http request", zap.Error(err))
		return utils.ErrSomethingWentWrong
	}
	defer resp.Body.Close()

	var body utils.Success[any]
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		c.log.Error("failed to decode response body", zap.Error(err))
		return utils.ErrSomethingWentWrong
	}

	if !body.Success {
		c.log.Error(
			"shipment service returned error",
			zap.Int("status", body.Code),
			zap.String("message", body.Message),
		)
		return &utils.HTTPError{
			Status:  resp.StatusCode,
			Message: body.Message,
		}
	}

	return nil
}

func (c *ShipmentClient) GetShipmentByTrackingID(ctx context.Context, trackingID uuid.UUID) (dto.ShipmentResponse, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/v1/shipment?trackingID=%s", c.shipmentSvcURL, trackingID.String()), nil)
	if err != nil {
		c.log.Error("failed to create http request", zap.Error(err))
		return dto.ShipmentResponse{}, utils.ErrSomethingWentWrong
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Error("failed to send http request", zap.Error(err))
		return dto.ShipmentResponse{}, utils.ErrSomethingWentWrong
	}
	defer resp.Body.Close()

	var body utils.Success[dto.ShipmentResponse]
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		c.log.Error("failed to decode response body", zap.Error(err))
		return dto.ShipmentResponse{}, utils.ErrSomethingWentWrong
	}

	if !body.Success {
		c.log.Error(
			"shipment service returned error",
			zap.Int("status", body.Code),
			zap.String("message", body.Message),
		)
		return dto.ShipmentResponse{}, &utils.HTTPError{
			Status:  resp.StatusCode,
			Message: body.Message,
		}
	}

	return body.Data, nil
}
