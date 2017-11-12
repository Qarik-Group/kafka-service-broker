package brokerconfig

import "os"

// Config contains the broker's primary configuration
type Config struct {
	KafkaConfiguration KafkaConfiguration
}

// KafkaConfiguration contains location/credentials for Kafka
type KafkaConfiguration struct {
	KafkaHostnames string
}

// LoadConfig loads environment variables into Config
func LoadConfig() (config Config) {
	config.KafkaConfiguration.KafkaHostnames = os.Getenv("KAFKA_HOSTNAMES")
	if config.KafkaConfiguration.KafkaHostnames == "" {
		config.KafkaConfiguration.KafkaHostnames = "localhost:9092"
	}
	return
}
