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

// SharedPlanRepository describes the creation/binding of a shared kafka service instances
// Unlike TopicPlanRepository, SharedPlanRepository does not provide any pre-created topics.
// Instead it is assumed that an orchestrator or producer will dynamically create topics and
// have their own way to share the topics with consumers.
// They should use the "topicNamePrefix" as the prefix for all topic names.
// Like TopicPlanRepository, Deprovision will delete all Kafka topics that have the service instanceID
// as a prefix on the topic name.
//
// Note, SharedPlanRepository currently still does create an initial topic (with the name instanceID)
// so as to indicate that the service instance has been already provisioned. End users can use it if they like.
type SharedPlanRepository struct {
	kafkaConfig brokerconfig.KafkaConfiguration
	logger      lager.Logger
}

// NewSharedPlanRepository creates a SharedPlanRepository
func NewSharedPlanRepository(kafkaConfig brokerconfig.KafkaConfiguration, logger lager.Logger) *SharedPlanRepository {
	return &SharedPlanRepository{
		kafkaConfig: kafkaConfig,
		logger:      logger,
	}
}

// InstanceExists returns true if instanceID belongs to an existing service instance
func (repo *SharedPlanRepository) InstanceExists(instanceID string) (bool, error) {
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
func (repo *SharedPlanRepository) Create(instanceID string) error {
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
		"plan":        "shared",
		"message":     "Successfully provisioned Kafka shared plan instance",
	})

	return nil
}

// Destroy will destroy any topics associated with the service instance
// Currently "associated with" is inferred - any topic name with instanceID as a prefix
func (repo *SharedPlanRepository) Destroy(instanceID string) error {
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
						"plan":        "shared",
						"topic.name":  topic.Name,
						"message":     "Failed to delete Kafka topic",
					})
				} else {
					repo.logger.Info("deprovision-instance.delete-topic", lager.Data{
						"instance_id": instanceID,
						"plan":        "shared",
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
		"plan":        "shared",
		"message":     "Successfully deprovisioned Kafka shared plan instance",
	})

	return nil
}

// Bind provides the credentials to access the Kafka cluster and the provided topics
func (repo *SharedPlanRepository) Bind(instanceID string, bindingID string) (broker.InstanceCredentials, error) {
	repo.logger.Info("bind-instance", lager.Data{
		"instance_id": instanceID,
		"binding_id":  bindingID,
		"plan":        "shared",
		"message":     "Successful bind of Kafka shared plan instance",
	})
	return broker.InstanceCredentials{
		ZookeeperPeers:  repo.kafkaConfig.ZookeeperPeers,
		KafkaHostnames:  repo.kafkaConfig.KafkaHostnames,
		TopicNamePrefix: instanceID,
	}, nil
}

// Unbind is a no-op as bindings are shared across all instances
func (repo *SharedPlanRepository) Unbind(instanceID string, bindingID string) error {
	repo.logger.Info("unbind-instance", lager.Data{
		"instance_id": instanceID,
		"binding_id":  bindingID,
		"plan":        "shared",
		"message":     "Successful unbind of Kafka shared plan instance",
	})
	return nil
}
