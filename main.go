package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"kafka-consumer/internal/config"
	"kafka-consumer/internal/consumer"
	"kafka-consumer/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	kafkaConsumer, err := consumer.NewKafkaConsumer(cfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer kafkaConsumer.Close()
	processor := service.NewProcessor(cfg)

	go func() {
		if err := kafkaConsumer.Consume(processor); err != nil {
			log.Printf("Error consuming messages: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gracefully...")
} 