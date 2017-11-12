package brokerconfig

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
)

// Config contains the broker's primary configuration
type Config struct {
	KafkaConfiguration KafkaConfiguration
	RedisConfiguration RedisConfiguration
}

// KafkaConfiguration contains location/credentials for Kafka
type KafkaConfiguration struct {
	KafkaHostnames string
}

// RedisConfiguration contains location/credentials for Redis used for internal storage
type RedisConfiguration struct {
	URI string
}

// LoadConfig loads environment variables into Config
func LoadConfig() (config Config, err error) {
	config.KafkaConfiguration.KafkaHostnames = os.Getenv("KAFKA_HOSTNAMES")
	if config.KafkaConfiguration.KafkaHostnames == "" {
		config.KafkaConfiguration.KafkaHostnames = "localhost:9092"
	}
	if os.Getenv("REDIS_URI") != "" {
		config.RedisConfiguration.URI = os.Getenv("REDIS_URI")
	} else if cfenv.IsRunningOnCF() {
		appEnv, _ := cfenv.Current()
		redisService, err := appEnv.Services.WithTag("redis")
		if err != nil {
			return config, err
		}
		if redisService[0].Credentials["uri"] == nil {
			return config, fmt.Errorf("expected redis instance %s to contain credentials 'uri'", redisService[0].Name)
		}
		config.RedisConfiguration.URI = redisService[0].Credentials["uri"].(string)
	}
	if config.RedisConfiguration.URI == "" {
		config.RedisConfiguration.URI = "redis://localhost:6379"
	}
	return
}
