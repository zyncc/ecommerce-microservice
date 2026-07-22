# 🛒 E-Commerce Microservices Platform

A production-minded e-commerce backend written in **Go**, built with a **microservices architecture** and an **event-driven design**. The platform exposes synchronous functionality through REST APIs, while asynchronous business workflows are coordinated through **Apache Kafka**.

It is designed to demonstrate the systems concerns behind a modern backend: independent services, resilient event processing, caching, authentication, rate limiting, container orchestration, CI, and observability.

🔐 The platform follows the **API Gateway pattern**: the API Gateway is the only public entry point for every API. Internal microservices remain private within the cluster, reducing the attack surface and centralizing authentication, authorization, rate limiting, and request handling at a single security boundary.

## ✨ Services at a Glance

| Service          | Responsibility                                                                                                                                 |
| ---------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| **API Gateway**  | Public entry point that routes requests to internal services, applies JWT authorization, exposes Swagger, and enforces IP-based rate limiting. |
| **Auth**         | Manages user registration, sign-in, JWT access/refresh token lifecycle, sessions, and customer addresses.                                      |
| **Product**      | Owns the product catalog and product-management operations.                                                                                    |
| **Inventory**    | Tracks product stock and reacts to successful-payment events to update inventory.                                                              |
| **Order**        | Creates and retrieves orders; consumes payment and shipment events to keep order state current.                                                |
| **Payment**      | Processes payment webhooks, persists payment records, and publishes `payment.succeeded` events.                                                |
| **Shipping**     | Creates shipments after successful payments, handles shipment webhooks, and publishes `shipment.updated` events.                               |
| **Notification** | Consumes successful-payment events and drives customer-notification workflows.                                                                 |

## 🚀 Run the Full Stack Locally

### Prerequisites

The full application is deployed to a local Kubernetes cluster. Install all of the following before continuing:

- 🐳 [Docker](https://www.docker.com/)
- ☸️ A local Kubernetes provider — this project uses [kind](https://kind.sigs.k8s.io/)
- `kubectl`
- [Helm](https://helm.sh/)

### 1. Create the kind cluster

From the repository root, create the cluster using the supplied port mappings for ingress:

```bash
kind create cluster --config kind-config.yaml
```

### 2. Install NGINX Ingress Controller

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```

Wait for the controller to become ready:

```bash
kubectl wait \
  --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s
```

### 3. Deploy the platform with Helm

This installs the microservices together with PostgreSQL, Redis, Kafka, and database migrations into the `ecommerce-microservice` namespace.

```bash
helm install ecommerce ./helm \
  --namespace ecommerce-microservice \
  --create-namespace
```

### 4. Wait for the services to become healthy

Startup can take a couple of minutes while images are pulled and the application dependencies initialize. Watch the rollout with:

```bash
kubectl get pods --namespace ecommerce-microservice --watch
```

When the pods are running, explore the API through Swagger at [http://localhost/swagger](http://localhost/swagger). 📚

## 🏗️ Architecture

The API Gateway is the single externally exposed application entry point. It forwards synchronous REST calls to the service that owns the relevant domain, keeping the public API separate from internal service topology.

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

The shipping service exposes a single shipment webhook endpoint. A shipment starts in `WAITING_FOR_PICKUP`; a webhook is required to advance it through each remaining stage:

```text
WAITING_FOR_PICKUP → PICKUP_DONE → SHIPPED → OUT_FOR_DELIVERY → DELIVERED
                       ▲              ▲             ▲                ▲
                       └──── Shipment webhook required for every transition ────┘
```

## 🧰 Technology Stack

| Area                      | Technologies                                           | How they are used                                                                                                   |
| ------------------------- | ------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------- |
| **Language & HTTP**       | Go, `net/http`, Chi                                    | Go services with lightweight, explicit REST routing and middleware.                                                 |
| **Architecture**          | Microservices, REST, event-driven messaging            | REST handles synchronous request/response work; Kafka handles asynchronous domain events.                           |
| **Messaging**             | Apache Kafka, Sarama                                   | Durable event transport and consumer groups for payment and shipment workflows.                                     |
| **Data**                  | PostgreSQL, pgx, SQL migrations                        | Relational persistence for users, products, inventory, orders, payments, and shipments.                             |
| **Caching & protection**  | Redis                                                  | Caches frequently accessed data and backs the IP-based token-bucket rate limiter at the API Gateway.                |
| **Security**              | JWT                                                    | Access and refresh tokens secure authenticated endpoints and session flows.                                         |
| **API documentation**     | Swagger / OpenAPI                                      | Interactive API documentation served by the gateway at `/swagger`.                                                  |
| **Containers & delivery** | Docker, GitHub Actions, GitHub Container Registry      | Each service is containerized; CI builds and publishes the latest images when code is pushed to `main`.             |
| **Orchestration**         | Kubernetes, Helm, kind, NGINX Ingress                  | Helm packages the full deployment; kind provides a reproducible local cluster; ingress exposes the gateway.         |
| **Observability**         | Structured logging, Prometheus, Grafana, Grafana Alloy | Production-ready service logs and monitoring for CPU/RAM usage, p99 request latency, and other operational signals. |

## 🔐 Security and Traffic Control

- **JWT authentication:** sign-in issues access and refresh tokens, and the gateway protects authenticated and administrator-only routes.
- **IP-based token bucket:** Redis-backed gateway middleware limits request bursts and helps protect the public API.
- **Idempotent integrations:** payment and shipping webhook flows use idempotency keys to safely handle retries.
- **Dead-letter topics:** poisoned messages and non-transient processing failures are isolated for investigation and controlled recovery.
- **HTTP resilience:** the gateway applies request IDs, structured request logging, panic recovery, CORS handling, and bounded downstream HTTP timeouts.

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
