#!/usr/bin/env bash

set -euo pipefail

echo "Installing NGINX Ingress Controller..."

helm upgrade --install ingress-nginx \
    ingress-nginx/ingress-nginx \
    -n ingress-nginx \
    --create-namespace

echo "Waiting for NGINX Ingress Controller..."

kubectl rollout status deployment/ingress-nginx-controller \
    -n ingress-nginx \
    --timeout=300s

echo "Installing Confluent for Kubernetes Operator..."

helm upgrade --install kafka \
    confluentinc/confluent-for-kubernetes \
    -n confluent \
    --create-namespace

echo "Waiting for Confluent Operator..."

kubectl rollout status deployment/confluent-operator \
    -n confluent \
    --timeout=300s

echo "Deploying Kafka Cluster..."

kubectl apply -R -f helm/kafka

echo "Waiting for KRaft Controller..."

kubectl rollout status \
    statefulset/kraftcontroller \
    -n confluent \
    --timeout=300s

echo "Waiting for Kafka broker..."

kubectl rollout status \
    statefulset/kafka \
    -n confluent \
    --timeout=300s

echo "Installing PostgreSQL..."

helm upgrade --install postgres \
    oci://registry-1.docker.io/bitnamicharts/postgresql \
    -n db \
    --create-namespace \
    -f helm/postgres/postgres-values.yaml

echo "Waiting for PostgreSQL to become ready..."

kubectl wait \
    --for=condition=ready pod \
    -l app.kubernetes.io/name=postgresql \
    -n db \
    --timeout=300s

echo "Running database migrations..."

helm upgrade --install migrations \
    ./helm/migrations \
    -n ecommerce \
    --create-namespace

echo "Deploying microservices..."

helm upgrade --install microservices \
    ./helm/microservices \
    -n ecommerce \
    --create-namespace

echo "Waiting for microservices..."

kubectl rollout status deployment/api-gateway -n ecommerce --timeout=300s
kubectl rollout status deployment/auth-service -n ecommerce --timeout=300s
kubectl rollout status deployment/product-service -n ecommerce --timeout=300s
kubectl rollout status deployment/order-service -n ecommerce --timeout=300s
kubectl rollout status deployment/payment-service -n ecommerce --timeout=300s
kubectl rollout status deployment/inventory-service -n ecommerce --timeout=300s
kubectl rollout status deployment/shipping-service -n ecommerce --timeout=300s

echo "Deploying Monitoring..."

helm upgrade --install prometheus-stack \
    oci://ghcr.io/prometheus-community/charts/kube-prometheus-stack \
    -n monitoring \
    --create-namespace

kubectl rollout status deployment/prometheus-stack-operator \
    -n monitoring \
    --timeout=300s

helm upgrade --install alloy \
    grafana/alloy \
    -n monitoring \
    --create-namespace \
    -f ./monitoring/alloy.yaml

kubectl rollout status deployment/alloy \
    -n monitoring \
    --timeout=300s

helm upgrade --install loki \
    grafana/loki \
    -n monitoring \
    --create-namespace \
    -f ./monitoring/loki.yaml

kubectl rollout status statefulset/loki \
    -n monitoring \
    --timeout=300s

echo "🚀 Deployment complete! 🚀"