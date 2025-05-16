# Build stage
FROM golang:1.24.1-bullseye AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -v -o kafka-consumer

FROM debian:bullseye-slim

WORKDIR /app

RUN apt-get update && apt-get install -y \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/kafka-consumer .
COPY --from=builder /app/.env .

CMD ["./kafka-consumer"] 