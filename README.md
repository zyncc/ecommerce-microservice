# 🛒 E-Commerce Microservices Platform

A production-minded e-commerce backend written in **Go**, built with a **microservices architecture** and an **event-driven design**. The API Gateway exposes a public REST/HTTP API, translates synchronous requests to **gRPC** calls for internal services, and coordinates asynchronous business workflows through **Apache Kafka**.

> **📈 Built to scale event processing:** Each event-consuming microservice runs as its own Kafka consumer group, so every service receives the events it needs while its replicas share the work within that service. High-throughput topics are split across multiple partitions; paired with Kubernetes replicas, Kafka assigns partitions across consumer instances to process events concurrently. This enables horizontal scaling and efficient, partition-level parallelism without duplicate processing within a consumer group.

It is designed to demonstrate the systems concerns behind a modern backend: independent services, resilient event processing, caching, authentication, rate limiting, container orchestration, CI, and observability.

🔐 The platform follows the **API Gateway pattern**: the API Gateway is the only public entry point for every API. Internal microservices remain private within the cluster, reducing the attack surface and centralizing authentication, authorization, rate limiting, and request handling at a single security boundary.

## ✨ Services at a Glance

| Service          | Responsibility                                                                                                                                                     |
| ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **API Gateway**  | Public REST/HTTP entry point that routes requests to internal services over gRPC, applies JWT authorization, exposes Swagger, and enforces IP-based rate limiting. |
| **Auth**         | Manages user registration, sign-in, JWT access/refresh token lifecycle, sessions, and customer addresses.                                                          |
| **Product**      | Owns the product catalog and product-management operations.                                                                                                        |
| **Inventory**    | Tracks product stock and reacts to successful-payment events to update inventory.                                                                                  |
| **Order**        | Creates and retrieves orders; consumes payment and shipment events to keep order state current.                                                                    |
| **Payment**      | Processes payment webhooks, persists payment records, and publishes `payment.succeeded` events.                                                                    |
| **Shipping**     | Creates shipments after successful payments, handles shipment webhooks, and publishes `shipment.updated` events.                                                   |
| **Notification** | Consumes successful-payment events and drives customer-notification workflows.                                                                                     |

## 🚀 Deploy the Full Stack

With `kubectl` configured for the target Kubernetes cluster, run the deployment script from the repository root:

```bash
./deploy.sh
```

The script installs the ingress controller, Kafka operator and cluster, PostgreSQL, database migrations, microservices, and monitoring stack, waiting for the required dependencies as it proceeds.

When running the cluster locally with kind, start `cloud-provider-kind` after the deployment. The NGINX Ingress Controller uses a `LoadBalancer` service, and kind does not assign its external IP automatically. `cloud-provider-kind` runs as a Docker container and assigns external IP addresses to the ingress controller and services such as Grafana:

```bash
docker run --network kind \
  registry.k8s.io/cloud-provider-kind/cloud-controller-manager
```

Keep this command running in a separate terminal. Then retrieve the assigned ingress address:

```bash
kubectl get ingress -A
```

Use the value in the `ADDRESS` column in place of `localhost` to access Swagger, Grafana, and the other APIs. 📚

## 🏗️ Architecture

The API Gateway is the single externally exposed application entry point. It accepts synchronous REST/HTTP requests and forwards them to the service that owns the relevant domain over gRPC, keeping the public API separate from internal service topology. Domain services that expose synchronous APIs do so through gRPC.

For long-running or cross-domain work, services communicate asynchronously through Kafka. A successful payment, for example, publishes a `payment.succeeded` event. Inventory, order, shipping, and notification services consume that event independently: stock can be adjusted, the order updated, a shipment created, and customer communication initiated without coupling these operations to the payment request. Shipping later emits `shipment.updated`, which the order service consumes to reflect delivery progress.

Poisoned messages and non-transient processing failures are routed to **dead-letter topics (DLTs)**. This prevents permanently invalid events from blocking consumers while preserving them for inspection, remediation, and controlled replay.

## 🛍️ Customer Order Journey

The platform models the full customer path from account creation to delivery. An order is created with a saved address, but it is not valid until payment succeeds. The payment webhook simulates provider confirmation: it marks the order as **PAID** and starts the asynchronous fulfillment flow.

```text
┌──────────┐    sign up / sign in    ┌──────────────┐
│ Customer │ ───────────────────────▶│ Auth Service │
└──────────┘                         └──────┬───────┘
     │                                      │ JWT access + refresh tokens
     │ create address                       ▼
     ├──────────────────────────────▶ Address saved
     │
     │ place order with saved address
     ▼
┌───────────────┐      pending payment      ┌─────────────────┐
│ Order Service │ ◀──────────────────────── │ API Gateway     │
└───────┬───────┘                           └─────────────────┘
        │
        │ payment-provider callback (mock)
        ▼
┌──────────────────┐  payment.succeeded  ┌─────────────────────────────┐
│ Payment Webhook  │ ───────────────────▶ │ Kafka                       │
│ → order is PAID  │                      │ event-driven fulfillment    │
└──────────────────┘                      └───────┬─────────┬──────────┘
                                                    │         │
                                      ┌─────────────┘         └─────────────┐
                                      ▼                                     ▼
                           ┌───────────────────┐                  ┌──────────────────────┐
                           │ Inventory Service │                  │ Notification Service │
                           │ update stock      │                  │ send customer email  │
                           └───────────────────┘                  └──────────────────────┘
                                      │
                                      ▼
                           ┌───────────────────┐  shipment.updated  ┌───────────────┐
                           │ Shipping Service  │ ──────────────────▶ │ Order Service │
                           │ create shipment   │                     │ update status │
                           └───────────────────┘                     └───────────────┘
```

### Shipment lifecycle

The API Gateway exposes a single shipment webhook endpoint, which invokes the shipping service over gRPC. A shipment starts in `WAITING_FOR_PICKUP`; a webhook is required to advance it through each remaining stage:

```text
WAITING_FOR_PICKUP → PICKUP_DONE → SHIPPED → OUT_FOR_DELIVERY → DELIVERED
                       ▲              ▲             ▲                ▲
                       └──── Shipment webhook required for every transition ────┘
```

## 🧰 Technology Stack

| Area                      | Technologies                                              | How they are used                                                                                                        |
| ------------------------- | --------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| **Language & transport**  | Go, `net/http`, Chi, gRPC                                 | Chi and `net/http` power the public gateway; internal services expose gRPC APIs.                                         |
| **Architecture**          | Microservices, REST gateway, gRPC, event-driven messaging | The gateway handles public REST/HTTP requests, gRPC handles synchronous internal calls, and Kafka carries domain events. |
| **Messaging**             | Apache Kafka, Sarama                                      | Durable event transport and consumer groups for payment and shipment workflows.                                          |
| **Data**                  | PostgreSQL, pgx, SQL migrations                           | Relational persistence for users, products, inventory, orders, payments, and shipments.                                  |
| **Caching & protection**  | Redis                                                     | Caches frequently accessed data and backs the IP-based token-bucket rate limiter at the API Gateway.                     |
| **Security**              | JWT                                                       | Access and refresh tokens secure authenticated endpoints and session flows.                                              |
| **API documentation**     | Swagger / OpenAPI                                         | Interactive API documentation served by the gateway at `/swagger`.                                                       |
| **Containers & delivery** | Docker, GitHub Actions, GitHub Container Registry         | Each service is containerized; CI builds and publishes the latest images when code is pushed to `main`.                  |
| **Orchestration**         | Kubernetes, Helm, kind, NGINX Ingress                     | Helm packages the full deployment; kind provides a reproducible local cluster; ingress exposes the gateway.              |
| **Observability**         | Structured logging, Prometheus, Grafana, Grafana Alloy    | Production-ready service logs and monitoring for CPU/RAM usage, p99 request latency, and other operational signals.      |

## 🔐 Security and Traffic Control

- **JWT authentication:** sign-in issues access and refresh tokens, and the gateway protects authenticated and administrator-only routes.
- **IP-based token bucket:** Redis-backed gateway middleware limits request bursts and helps protect the public API.
- **Idempotent integrations:** payment and shipping webhook flows use idempotency keys to safely handle retries.
- **Dead-letter topics:** poisoned messages and non-transient processing failures are isolated for investigation and controlled recovery.
- **Gateway resilience:** the gateway applies request IDs, structured request logging, panic recovery, and CORS handling while brokering public HTTP requests to gRPC-backed services.

## 📈 Observability

Every microservice is instrumented for production-oriented logging, making it easier to trace requests and investigate failures across service boundaries. The observability stack uses:

- **Grafana Alloy** to collect and enrich container logs.
- **Prometheus** to collect infrastructure and application metrics.
- **Grafana** to visualize server resource consumption (CPU and RAM), request latency—including p99—and other health indicators.

Together, these tools make the platform observable beyond application correctness: they provide the operational feedback needed to understand performance and reliability under real traffic. 🔎

## 🔄 Continuous Integration

GitHub Actions builds Docker images for every microservice and the migration image. On pushes to `main`, the workflow tags and publishes the latest images to GitHub Container Registry, keeping deployable artifacts aligned with the main branch.

---

Built to showcase practical backend engineering across distributed systems, cloud-native deployment, and operational visibility. 🚀
