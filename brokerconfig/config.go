package brokerconfig

import (
	"os"
	"time"

	"github.com/wvanbergen/kazoo-go"
)

// Config contains the broker's primary configuration
type Config struct {
	Broker             BrokerConfiguration
	KafkaConfiguration KafkaConfiguration
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
	return
}
