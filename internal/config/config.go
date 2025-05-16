package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	KafkaBootstrapServers string
	KafkaGroupID          string
	KafkaTopic            string
	KafkaTopicDLQ         string
	TargetServiceURL      string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
	viper.SetDefault("KAFKA_GROUP_ID", "order-processor-group")
	viper.SetDefault("KAFKA_TOPIC", "order-status")
	viper.SetDefault("KAFKA_TOPIC_DLQ", "order-status-dlq")
	viper.SetDefault("TARGET_SERVICE_URL", "http://localhost:8080/api/v1")

	if err := viper.ReadInConfig(); err != nil {
		// If .env file doesn't exist, continue with default values
		fmt.Printf("Warning: .env file not found, using default values: %v\n", err)
	}

	config := &Config{
		KafkaBootstrapServers: viper.GetString("KAFKA_BOOTSTRAP_SERVERS"),
		KafkaGroupID:          viper.GetString("KAFKA_GROUP_ID"),
		KafkaTopic:            viper.GetString("KAFKA_TOPIC"),
		KafkaTopicDLQ:         viper.GetString("KAFKA_TOPIC_DLQ"),
		TargetServiceURL:      viper.GetString("TARGET_SERVICE_URL"),
	}

	return config, nil
}
