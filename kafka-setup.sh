#!/bin/bash

echo "Waiting for Kafka to be ready..."
sleep 30

echo "Creating topics..."
kafka-topics --bootstrap-server kafka:29092 --create --topic order-status --partitions 1 --replication-factor 1
kafka-topics --bootstrap-server kafka:29092 --create --topic order-status-dlq --partitions 1 --replication-factor 1

echo "Listing topics..."
kafka-topics --bootstrap-server kafka:29092 --list

echo "Kafka setup completed!" 