package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "localhost:9092", cfg.KafkaBootstrapServers)
	assert.Equal(t, "order-processor-group", cfg.KafkaGroupID)
	assert.Equal(t, "order-status", cfg.KafkaTopic)
	assert.Equal(t, "order-status-dlq", cfg.KafkaTopicDLQ)
	assert.Equal(t, "http://localhost:8080/api/v1", cfg.TargetServiceURL)

	os.Setenv("KAFKA_BOOTSTRAP_SERVERS", "kafka:29092")
	os.Setenv("KAFKA_GROUP_ID", "test-group")
	os.Setenv("KAFKA_TOPIC", "test-topic")
	os.Setenv("KAFKA_TOPIC_DLQ", "test-topic-dlq")
	os.Setenv("TARGET_SERVICE_URL", "http://test:8080/api/v1")

	cfg, err = Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "kafka:29092", cfg.KafkaBootstrapServers)
	assert.Equal(t, "test-group", cfg.KafkaGroupID)
	assert.Equal(t, "test-topic", cfg.KafkaTopic)
	assert.Equal(t, "test-topic-dlq", cfg.KafkaTopicDLQ)
	assert.Equal(t, "http://test:8080/api/v1", cfg.TargetServiceURL)

	os.Unsetenv("KAFKA_BOOTSTRAP_SERVERS")
	os.Unsetenv("KAFKA_GROUP_ID")
	os.Unsetenv("KAFKA_TOPIC")
	os.Unsetenv("KAFKA_TOPIC_DLQ")
	os.Unsetenv("TARGET_SERVICE_URL")
}

func TestLoad_WithEnvFile(t *testing.T) {
	envContent := `
KAFKA_BOOTSTRAP_SERVERS=kafka:29092
KAFKA_GROUP_ID=test-group
KAFKA_TOPIC=test-topic
KAFKA_TOPIC_DLQ=test-topic-dlq
TARGET_SERVICE_URL=http://test:8080/api/v1
`
	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(".env")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "kafka:29092", cfg.KafkaBootstrapServers)
	assert.Equal(t, "test-group", cfg.KafkaGroupID)
	assert.Equal(t, "test-topic", cfg.KafkaTopic)
	assert.Equal(t, "test-topic-dlq", cfg.KafkaTopicDLQ)
	assert.Equal(t, "http://test:8080/api/v1", cfg.TargetServiceURL)
}
