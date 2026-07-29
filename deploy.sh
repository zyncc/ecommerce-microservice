#!/usr/bin/env bash

set -euo pipefail

echo "Installing Confluent for Kubernetes Operator..."

helm upgrade --install kafka \
    confluentinc/confluent-for-kubernetes \
    -n confluent \
    --create-namespace \
    --timeout=10m

echo "Waiting for Confluent Operator..."

kubectl rollout status deployment/confluent-operator \
    -n confluent \
    --timeout=300s

echo "Deploying Kafka Cluster..."

kubectl apply -R -f helm/kafka

echo "Installing NGINX Ingress Controller..."

helm upgrade --install ingress-nginx \
    ingress-nginx/ingress-nginx \
    -n ingress-nginx \
    --create-namespace \
    --wait \
    --timeout=10m

echo "Installing PostgreSQL..."

helm upgrade --install postgres \
    oci://registry-1.docker.io/bitnamicharts/postgresql \
    -n db \
    --create-namespace \
    -f helm/postgres/postgres-values.yaml \
    --wait \
    --timeout=10m

echo "Waiting for PostgreSQL to become ready..."

kubectl wait \
    --for=condition=Ready \
    pod \
    -l app.kubernetes.io/name=postgresql \
    -n db \
    --timeout=300s

echo "Running database migrations..."

helm upgrade --install migrations \
    ./helm/migrations \
    -n ecommerce \
    --create-namespace \
    --timeout=10m

echo "Waiting for KRaft Controller..."

kubectl rollout status \
    statefulset/kraftcontroller \
    -n confluent \
    --timeout=300s

echo "Waiting for Kafka StatefulSet to be created..."

until kubectl get statefulset kafka -n confluent >/dev/null 2>&1; do
    sleep 2
done

echo "Waiting for Kafka broker..."

kubectl rollout status \
    statefulset/kafka \
    -n confluent \
    --timeout=300s

echo "Waiting for Kafka topics..."

until kubectl exec -n confluent kafka-0 -- \
    kafka-topics --bootstrap-server localhost:9092 --list | \
    grep -q "payment.succeeded"; do
    sleep 2
done

until kubectl exec -n confluent kafka-0 -- \
    kafka-topics --bootstrap-server localhost:9092 --list | \
    grep -q "shipment.updated"; do
    sleep 2
done

echo "Deploying microservices..."

helm upgrade --install microservices \
    ./helm/microservices \
    -n ecommerce \
    --create-namespace \
    --timeout=10m

echo "Deploying Monitoring..."

helm upgrade --install prometheus-stack \
    oci://ghcr.io/prometheus-community/charts/kube-prometheus-stack \
    -n monitoring \
    -f monitoring/prometheus.yaml \
    --create-namespace \
    --timeout=10m

echo "Deploying Loki..."

helm upgrade --install loki \
    grafana/loki \
    -n monitoring \
    --create-namespace \
    -f ./monitoring/loki.yaml \
    --timeout=10m

echo "Deploying Alloy..."

helm upgrade --install alloy \
    grafana/alloy \
    -n monitoring \
    --create-namespace \
    -f ./monitoring/alloy.yaml \
    --timeout=10m

echo "Creating monitoring ingress..."

kubectl apply -f monitoring/ingress.yaml

echo "🚀 Deployment complete! 🚀"