package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	"kafka-consumer/internal/config"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Payload struct {
	OrdemDeVenda string `json:"ordemDeVenda"`
	EtapaAtual   string `json:"etapaAtual"`
}

type Processor interface {
	Process(payload *Payload) error
}

type KafkaConsumer struct {
	consumer *kafka.Consumer
	topic    string
	dlqTopic string
	config   *config.Config
}

func NewKafkaConsumer(cfg *config.Config) (*KafkaConsumer, error) {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":  cfg.KafkaBootstrapServers,
		"group.id":          cfg.KafkaGroupID,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &KafkaConsumer{
		consumer: consumer,
		topic:    cfg.KafkaTopic,
		dlqTopic: cfg.KafkaTopicDLQ,
		config:   cfg,
	}, nil
}

func (k *KafkaConsumer) Consume(processor Processor) error {
	err := k.consumer.Subscribe(k.topic, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	for {
		msg, err := k.consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var payload Payload
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			if err := k.sendToDLQ(msg.Value); err != nil {
				log.Printf("Error sending to DLQ: %v", err)
			}
			continue
		}

		if err := processor.Process(&payload); err != nil {
			log.Printf("Error processing message: %v", err)
			if err := k.sendToDLQ(msg.Value); err != nil {
				log.Printf("Error sending to DLQ: %v", err)
			}
			continue
		}

		if _, err := k.consumer.CommitMessage(msg); err != nil {
			log.Printf("Error committing message: %v", err)
		}
	}
}

func (k *KafkaConsumer) sendToDLQ(message []byte) error {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": k.config.KafkaBootstrapServers,
	}

	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		return fmt.Errorf("failed to create producer for DLQ: %w", err)
	}
	defer producer.Close()

	deliveryChan := make(chan kafka.Event)
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &k.dlqTopic,
			Partition: kafka.PartitionAny,
		},
		Value: message,
	}, deliveryChan)

	if err != nil {
		return fmt.Errorf("failed to produce message to DLQ: %w", err)
	}

	<-deliveryChan
	return nil
}

func (k *KafkaConsumer) Close() {
	if k.consumer != nil {
		k.consumer.Close()
	}
} 