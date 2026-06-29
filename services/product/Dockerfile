FROM golang:1.26.4-trixie AS builder

WORKDIR /app
COPY go.work go.work.sum ./

COPY services/api-gateway ./services/api-gateway
COPY services/auth ./services/auth
COPY services/product ./services/product

RUN go work sync

WORKDIR /app/services/product

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /product-service ./cmd/api

FROM alpine:3.23.5

# required if service makes external HTTPS calls
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /product-service /product-service
EXPOSE 8082

CMD [ "/product-service" ]
