package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	KafkaBootstrapServers string
	KafkaGroupID         string
	KafkaTopic           string
	KafkaTopicDLQ        string
	TargetServiceURL     string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{
		KafkaBootstrapServers: viper.GetString("KAFKA_BOOTSTRAP_SERVERS"),
		KafkaGroupID:         viper.GetString("KAFKA_GROUP_ID"),
		KafkaTopic:           viper.GetString("KAFKA_TOPIC"),
		KafkaTopicDLQ:        viper.GetString("KAFKA_TOPIC_DLQ"),
		TargetServiceURL:     viper.GetString("TARGET_SERVICE_URL"),
	}

	return config, nil
} 