#!/bin/bash

echo "Waiting for Kafka to be ready..."
sleep 30

# Create topics
kafka-topics --bootstrap-server kafka:29092 \
    --create --if-not-exists \
    --topic order-status \
    --partitions 1 \
    --replication-factor 1

kafka-topics --bootstrap-server kafka:29092 \
    --create --if-not-exists \
    --topic order-status-dlq \
    --partitions 1 \
    --replication-factor 1

# List topics to verify
kafka-topics --bootstrap-server kafka:29092 --list

echo "Kafka setup completed" 