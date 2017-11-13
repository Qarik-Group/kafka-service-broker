package brokerconfig

import (
	"fmt"
	"os"
	"time"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/wvanbergen/kazoo-go"
)

// Config contains the broker's primary configuration
type Config struct {
	Broker             BrokerConfiguration
	KafkaConfiguration KafkaConfiguration
	RedisConfiguration RedisConfiguration
}

// BrokerConfiguration contains the auth credentials
type BrokerConfiguration struct {
	ListenPort string
	Username   string
	Password   string
}

// KafkaConfiguration contains location/credentials for Kafka
type KafkaConfiguration struct {
	ZookeeperPeers         string
	ZookeeperTimeout       time.Duration
	KafkaHostnames         string
	KafkaPartitionCount    int
	KafkaReplicationFactor int
}

// RedisConfiguration contains location/credentials for Redis used for internal storage
type RedisConfiguration struct {
	URI string
}

// LoadConfig loads environment variables into Config
func LoadConfig() (config Config, err error) {
	config.Broker.ListenPort = os.Getenv("PORT")
	if config.Broker.ListenPort == "" {
		config.Broker.ListenPort = "8100"
	}
	config.Broker.Username = os.Getenv("BROKER_USERNAME")
	config.Broker.Password = os.Getenv("BROKER_PASSWORD")

	config.KafkaConfiguration.ZookeeperPeers = os.Getenv("ZOOKEEPER_PEERS")
	if config.KafkaConfiguration.ZookeeperPeers == "" {
		config.KafkaConfiguration.ZookeeperPeers = "localhost:2181"
	}
	config.KafkaConfiguration.ZookeeperTimeout = 1000

	zkConf := kazoo.NewConfig()
	zkConf.Timeout = time.Duration(config.KafkaConfiguration.ZookeeperTimeout) * time.Millisecond

	kz, err := kazoo.NewKazooFromConnectionString(config.KafkaConfiguration.ZookeeperPeers, zkConf)
	if err != nil {
		return
	}
	defer func() { _ = kz.Close() }()

	brokers, err := kz.BrokerList()
	if err != nil {
		return
	}

	for _, broker := range brokers {
		if config.KafkaConfiguration.KafkaHostnames != "" {
			config.KafkaConfiguration.KafkaHostnames += ","
		}
		config.KafkaConfiguration.KafkaHostnames += broker
	}
	config.KafkaConfiguration.KafkaReplicationFactor = len(brokers)
	config.KafkaConfiguration.KafkaPartitionCount = 2

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
