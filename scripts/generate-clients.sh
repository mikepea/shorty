#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

echo "==> Generating Swagger 2.0 docs from Go annotations..."
swag init -g cmd/shorty-server/main.go -o "$REPO_ROOT/api/swagger"

echo "==> Converting Swagger 2.0 → OpenAPI 3.0..."
npx --yes swagger2openapi "$REPO_ROOT/api/swagger/swagger.yaml" \
  -o "$REPO_ROOT/api/swagger/openapi3.yaml" --yaml

echo "==> Generating Go client from OpenAPI 3.0 spec..."
mkdir -p "$REPO_ROOT/clients/go/client"
oapi-codegen -package client \
  -generate types,client \
  "$REPO_ROOT/api/swagger/openapi3.yaml" \
  > "$REPO_ROOT/clients/go/client/client.gen.go"

echo "==> Client generation complete."
