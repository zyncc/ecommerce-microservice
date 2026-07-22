FROM golang:1.26.4-trixie AS builder

WORKDIR /app
COPY go.work go.work.sum ./

COPY services/api-gateway ./services/api-gateway
COPY services/auth ./services/auth
COPY services/product ./services/product
COPY services/inventory ./services/inventory
COPY services/order ./services/order
COPY services/payment ./services/payment
COPY services/shipping ./services/shipping
COPY services/notification ./services/notification

RUN go work sync

WORKDIR /app/services/api-gateway

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /api-gateway ./cmd/api

FROM alpine:3.23.5

# required if service makes external HTTPS calls
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /api-gateway /api-gateway
EXPOSE 8080

CMD [ "/api-gateway" ]
