package kafka

import (
	"code.cloudfoundry.org/lager"

	"github.com/starkandwayne/kafka-service-broker/broker"
	"github.com/starkandwayne/kafka-service-broker/brokerconfig"
)

// TopicRepository describes the creation/binding of topic-orientated kafka service instances
type TopicRepository struct {
	kafkaConfig       brokerconfig.KafkaConfiguration
	logger            lager.Logger
	inMemInstanceList []string
}

// NewTopicRepository creates a TopicRepository
func NewTopicRepository(kafkaConfig brokerconfig.KafkaConfiguration, logger lager.Logger) *TopicRepository {
	return &TopicRepository{
		kafkaConfig: kafkaConfig,
		logger:      logger,
	}
}

// InstanceExists returns true if instanceID belongs to an existing service instance
func (repo *TopicRepository) InstanceExists(instanceID string) (bool, error) {
	for _, inMem := range repo.inMemInstanceList {
		if inMem == instanceID {
			return true, nil
		}
	}
	return false, nil
}

// Create will create a topic(s)
func (repo *TopicRepository) Create(instanceID string) error {
	repo.inMemInstanceList = append(repo.inMemInstanceList, instanceID)
	repo.logger.Info("provision-instance", lager.Data{
		"instance_id": instanceID,
		"plan":        "topic",
		"message":     "Successfully provisioned Kafka topic instance",
	})

	return nil
}

// Destroy will destroy any topics associated with the service instance
func (repo *TopicRepository) Destroy(instanceID string) error {
	repo.logger.Info("deprovision-instance", lager.Data{
		"instance_id": instanceID,
		"plan":        "topic",
		"message":     "Successfully deprovisioned Kafka topic instance",
	})

	return nil
}

// Bind provides the credentials to access the Kafka cluster and the provided topics
func (repo *TopicRepository) Bind(instanceID string, bindingID string) (broker.InstanceCredentials, error) {
	repo.logger.Info("bind-instance", lager.Data{
		"instance_id": instanceID,
		"binding_id":  bindingID,
		"plan":        "topic",
		"message":     "Successful bind of Kafka topic instance",
	})
	return broker.InstanceCredentials{
		KafkaHostnames: repo.kafkaConfig.KafkaHostnames,
	}, nil
}

// Unbind is a no-op as bindings are shared across all instances
func (repo *TopicRepository) Unbind(instanceID string, bindingID string) error {
	repo.logger.Info("unbind-instance", lager.Data{
		"instance_id": instanceID,
		"binding_id":  bindingID,
		"plan":        "topic",
		"message":     "Successful unbind of Kafka topic instance",
	})
	return nil
}
