package kafka

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/wvanbergen/kazoo-go"

	"github.com/starkandwayne/kafka-service-broker/broker"
	"github.com/starkandwayne/kafka-service-broker/brokerconfig"
)

// TopicPlanRepository describes the creation/binding of topic-orientated kafka service instances
type TopicPlanRepository struct {
	kafkaConfig brokerconfig.KafkaConfiguration
	logger      lager.Logger
}

// NewTopicPlanRepository creates a TopicPlanRepository
func NewTopicPlanRepository(kafkaConfig brokerconfig.KafkaConfiguration, logger lager.Logger) *TopicPlanRepository {
	return &TopicPlanRepository{
		kafkaConfig: kafkaConfig,
		logger:      logger,
	}
}

// InstanceExists returns true if instanceID belongs to an existing service instance
func (repo *TopicPlanRepository) InstanceExists(instanceID string) (bool, error) {
	zkConf := kazoo.NewConfig()
	zkConf.Timeout = time.Duration(repo.kafkaConfig.ZookeeperTimeout) * time.Millisecond
	kz, err := kazoo.NewKazooFromConnectionString(repo.kafkaConfig.ZookeeperPeers, zkConf)
	if err != nil {
		return false, err
	}
	defer func() { _ = kz.Close() }()
	return kz.Topic(instanceID).Exists()
}

// Create will create a topic(s)
func (repo *TopicPlanRepository) Create(instanceID string) error {
	zkConf := kazoo.NewConfig()
	zkConf.Timeout = time.Duration(repo.kafkaConfig.ZookeeperTimeout) * time.Millisecond
	kz, err := kazoo.NewKazooFromConnectionString(repo.kafkaConfig.ZookeeperPeers, zkConf)
	if err != nil {
		return err
	}
	defer func() { _ = kz.Close() }()
	// A topic with the name of the instanceID is created, even if it is not returned
	// via credentials. It is currently used as proof that the service instance exists.
	kz.CreateTopic(instanceID,
		repo.kafkaConfig.KafkaPartitionCount,
		repo.kafkaConfig.KafkaReplicationFactor,
		map[string]string{})

	repo.logger.Info("provision-instance", lager.Data{
		"instance_id": instanceID,
		"plan":        "topic",
		"message":     "Successfully provisioned Kafka topic plan instance",
	})

	return nil
}

// Destroy will destroy any topics associated with the service instance
// Currently "associated with" is inferred - any topic name with instanceID as a prefix
func (repo *TopicPlanRepository) Destroy(instanceID string) error {
	zkConf := kazoo.NewConfig()
	zkConf.Timeout = time.Duration(repo.kafkaConfig.ZookeeperTimeout) * time.Millisecond
	kz, err := kazoo.NewKazooFromConnectionString(repo.kafkaConfig.ZookeeperPeers, zkConf)
	if err != nil {
		return err
	}
	defer func() { _ = kz.Close() }()
	allTopics, err := kz.Topics()
	if err != nil {
		return fmt.Errorf("Failed to get Kafka topics from Zookeeper: %v", err)
	}

	var wg sync.WaitGroup
	for i, topic := range allTopics {
		if strings.HasPrefix(topic.Name, instanceID) {
			wg.Add(1)
			go func(i int, topic *kazoo.Topic) {
				defer wg.Done()
				fmt.Println("Deleting topic", topic.Name)
				err = kz.DeleteTopic(topic.Name)
				if err != nil {
					repo.logger.Error("deprovision-instance.delete-topic", err, lager.Data{
						"instance_id": instanceID,
						"plan":        "topic",
						"topic.name":  topic.Name,
						"message":     "Failed to delete Kafka topic",
					})
				} else {
					repo.logger.Info("deprovision-instance.delete-topic", lager.Data{
						"instance_id": instanceID,
						"plan":        "topic",
						"topic.name":  topic.Name,
						"message":     "Successfully deleted Kafka topic",
					})
				}
			}(i, topic)
		}
	}

	wg.Wait()

	repo.logger.Info("deprovision-instance", lager.Data{
		"instance_id": instanceID,
		"plan":        "topic",
		"message":     "Successfully deprovisioned Kafka topic plan instance",
	})

	return nil
}

// Bind provides the credentials to access the Kafka cluster and the provided topics
func (repo *TopicPlanRepository) Bind(instanceID string, bindingID string) (broker.InstanceCredentials, error) {
	repo.logger.Info("bind-instance", lager.Data{
		"instance_id": instanceID,
		"binding_id":  bindingID,
		"plan":        "topic",
		"message":     "Successful bind of Kafka topic plan instance",
	})
	return broker.InstanceCredentials{
		ZookeeperPeers: repo.kafkaConfig.ZookeeperPeers,
		KafkaHostnames: repo.kafkaConfig.KafkaHostnames,
		TopicName:      instanceID,
	}, nil
}

// Unbind is a no-op as bindings are shared across all instances
func (repo *TopicPlanRepository) Unbind(instanceID string, bindingID string) error {
	repo.logger.Info("unbind-instance", lager.Data{
		"instance_id": instanceID,
		"binding_id":  bindingID,
		"plan":        "topic",
		"message":     "Successful unbind of Kafka topic plan instance",
	})
	return nil
}
