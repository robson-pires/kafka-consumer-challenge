#!/bin/bash

echo "Waiting for Kafka to be ready..."
sleep 30

TEST_MESSAGE='{"ordemDeVenda": "order12345", "etapaAtual": "FATURADO"}'

echo "Sending test message to Kafka..."
echo $TEST_MESSAGE | kafka-console-producer \
    --bootstrap-server kafka:29092 \
    --topic order-status

echo "Test message sent successfully!" 