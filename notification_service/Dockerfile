# Build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache \
    protoc \
    protobuf-dev \
    make \
    git \
    bash

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN chmod +x scripts/generate.sh
RUN ./scripts/generate.sh

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY config.yaml .

EXPOSE 8081 50052
CMD ["./main"]