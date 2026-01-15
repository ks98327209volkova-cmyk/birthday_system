# Build stage
FROM golang:1.24-alpine AS builder

# Устанавливаем protoc и необходимые инструменты для Alpine
RUN apk add --no-cache \
    protoc \
    protobuf-dev \
    make \
    git \
    bash

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Даем права на выполнение скрипта
RUN chmod +x scripts/generate.sh

# Генерируем protobuf код (внутри контейнера)
RUN ./scripts/generate.sh

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/main .

# Копируем конфиг
COPY config.yaml .

EXPOSE 8080 50051

CMD ["./main"]