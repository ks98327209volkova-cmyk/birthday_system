#!/bin/bash

cd "$(dirname "$0")/.." || exit

# Создаем директории если их нет
mkdir -p internal/pb/models
mkdir -p internal/pb/gift_manager_api
mkdir -p internal/pb/swagger/gift_manager_api

# Генерация gRPC кода
protoc -I ./api \
  -I ./api/google/api \
  --go_out=./internal/pb --go_opt=paths=source_relative \
  --go-grpc_out=./internal/pb --go-grpc_opt=paths=source_relative \
  ./api/gift_manager_api/gift_manager.proto ./api/models/person_model.proto

# Генерация gRPC-Gateway
protoc -I ./api \
  -I ./api/google/api \
  --grpc-gateway_out=./internal/pb \
  --grpc-gateway_opt paths=source_relative \
  --grpc-gateway_opt logtostderr=true \
  ./api/gift_manager_api/gift_manager.proto

# Генерация OpenAPI
protoc -I ./api \
  -I ./api/google/api \
  --openapiv2_out=./internal/pb/swagger \
  --openapiv2_opt logtostderr=true \
  ./api/gift_manager_api/gift_manager.proto

echo "Code generation for gift_manager complete!"