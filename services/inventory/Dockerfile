FROM golang:1.26.4-trixie AS builder

WORKDIR /app

COPY go.work .
COPY go.work.sum .

COPY services/api-gateway/go.mod services/api-gateway/go.sum ./services/api-gateway/
COPY services/auth/go.mod services/auth/go.sum ./services/auth/
COPY services/product/go.mod services/product/go.sum ./services/product/
COPY services/inventory/go.mod services/inventory/go.sum ./services/inventory/
COPY services/order/go.mod services/order/go.sum ./services/order/
COPY services/payment/go.mod services/payment/go.sum ./services/payment/
COPY services/shipping/go.mod services/shipping/go.sum ./services/shipping/
COPY services/notification/go.mod services/notification/go.sum ./services/notification/

RUN go work sync

COPY services ./services

WORKDIR /app/services/inventory

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -o /inventory-service ./cmd/api

FROM alpine:3.23.5

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /inventory-service /inventory-service

EXPOSE 8083

CMD ["/inventory-service"]