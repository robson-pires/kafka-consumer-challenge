#!/bin/bash

echo "Waiting for Kafka to be ready..."
sleep 30

# Define test message
TEST_MESSAGE='{"ordemDeVenda": "order12345", "etapaAtual": "FATURADO"}'

# Send test message
echo "Sending test message..."
echo "$TEST_MESSAGE" | kafka-console-producer \
    --bootstrap-server kafka:29092 \
    --topic order-status

echo "Test message sent successfully" 