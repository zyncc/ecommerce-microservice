#!/usr/bin/env bash

set -euo pipefail

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

echo "Deployment complete!"