package consumer

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"kafka-consumer/internal/config"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProcessor implements the Processor interface for testing
type MockProcessor struct {
	mock.Mock
}

func (m *MockProcessor) Process(payload *Payload) error {
	args := m.Called(payload)
	return args.Error(0)
}

// MockConsumer is a mock implementation of the Kafka consumer
type MockConsumer struct {
	mock.Mock
}

func (m *MockConsumer) Subscribe(topic string, rebalanceCb kafka.RebalanceCb) error {
	args := m.Called(topic, rebalanceCb)
	return args.Error(0)
}

func (m *MockConsumer) ReadMessage(timeout time.Duration) (*kafka.Message, error) {
	args := m.Called(timeout)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*kafka.Message), args.Error(1)
}

func (m *MockConsumer) CommitMessage(msg *kafka.Message) ([]kafka.TopicPartition, error) {
	args := m.Called(msg)
	return args.Get(0).([]kafka.TopicPartition), args.Error(1)
}

func (m *MockConsumer) Close() {
	m.Called()
}

func TestNewKafkaConsumer(t *testing.T) {
	cfg := &config.Config{
		KafkaBootstrapServers: "localhost:9092",
		KafkaGroupID:          "test-group",
		KafkaTopic:            "test-topic",
		KafkaTopicDLQ:         "test-topic-dlq",
	}

	consumer, err := NewKafkaConsumer(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, consumer)
	assert.Equal(t, cfg.KafkaTopic, consumer.topic)
	assert.Equal(t, cfg.KafkaTopicDLQ, consumer.dlqTopic)
}

func TestKafkaConsumer_Consume(t *testing.T) {
	tests := []struct {
		name          string
		message       *kafka.Message
		processError  error
		expectedError bool
		expectedDLQ   bool
	}{
		{
			name: "successful message processing",
			message: &kafka.Message{
				Value: []byte(`{"ordemDeVenda": "test123", "etapaAtual": "FATURADO"}`),
			},
			processError:  nil,
			expectedError: false,
			expectedDLQ:   false,
		},
		{
			name: "invalid message format",
			message: &kafka.Message{
				Value: []byte(`invalid json`),
			},
			processError:  nil,
			expectedError: true,
			expectedDLQ:   true,
		},
		{
			name: "processing error",
			message: &kafka.Message{
				Value: []byte(`{"ordemDeVenda": "test123", "etapaAtual": "FATURADO"}`),
			},
			processError:  assert.AnError,
			expectedError: true,
			expectedDLQ:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock processor
			mockProcessor := new(MockProcessor)
			if tt.message != nil {
				var payload Payload
				if err := json.Unmarshal(tt.message.Value, &payload); err == nil {
					mockProcessor.On("Process", &payload).Return(tt.processError)
				}
			}

			// Create mock consumer
			mockKafkaConsumer := new(MockConsumer)
			mockKafkaConsumer.On("Subscribe", mock.Anything, mock.Anything).Return(nil)
			mockKafkaConsumer.On("ReadMessage", mock.Anything).Return(tt.message, nil)
			if !tt.expectedError {
				mockKafkaConsumer.On("CommitMessage", mock.Anything).Return([]kafka.TopicPartition{}, nil)
			}

			// Create consumer with mock
			consumer := &KafkaConsumer{
				consumer: mockKafkaConsumer,
				topic:    "test-topic",
				dlqTopic: "test-topic-dlq",
				config:   &config.Config{},
			}

			// Create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			// Start consumer in goroutine
			errChan := make(chan error, 1)
			go func() {
				errChan <- consumer.Consume(mockProcessor)
			}()

			// Wait for context timeout or error
			select {
			case err := <-errChan:
				if tt.expectedError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			case <-ctx.Done():
				if !tt.expectedError {
					t.Error("Expected no error but got timeout")
				}
			}

			// Verify mock expectations
			mockProcessor.AssertExpectations(t)
			mockKafkaConsumer.AssertExpectations(t)
		})
	}
}

func TestKafkaConsumer_sendToDLQ(t *testing.T) {
	cfg := &config.Config{
		KafkaBootstrapServers: "localhost:9092",
		KafkaTopicDLQ:         "test-topic-dlq",
	}

	consumer := &KafkaConsumer{
		config:   cfg,
		dlqTopic: cfg.KafkaTopicDLQ,
	}

	message := []byte(`{"ordemDeVenda": "test123", "etapaAtual": "FATURADO"}`)
	err := consumer.sendToDLQ(message)
	assert.NoError(t, err)
}
