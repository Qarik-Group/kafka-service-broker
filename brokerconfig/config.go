package brokerconfig

import "os"

// Config contains the broker's primary configuration
type Config struct {
	KafkaHostnames string
}

// LoadConfig loads environment variables into Config
func LoadConfig() (config Config) {
	config.KafkaHostnames = os.Getenv("KAFKA_HOSTNAMES")
	if config.KafkaHostnames == "" {
		config.KafkaHostnames = "localhost:9092"
	}
	return
}
