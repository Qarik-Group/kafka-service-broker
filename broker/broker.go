package broker

import (
	"context"
	"errors"
	"fmt"

	"github.com/pivotal-cf/brokerapi"
	"github.com/starkandwayne/kafka-service-broker/brokerconfig"
)

type InstanceCredentials struct {
	ZookeeperPeers string
	KafkaHostnames string
	TopicName      string
}

type InstanceCreator interface {
	Create(instanceID string) error
	Destroy(instanceID string) error
	InstanceExists(instanceID string) (bool, error)
}

type InstanceBinder interface {
	Bind(instanceID string, bindingID string) (InstanceCredentials, error)
	Unbind(instanceID string, bindingID string) error
	InstanceExists(instanceID string) (bool, error)
}

type KafkaServiceBroker struct {
	InstanceCreators map[string]InstanceCreator
	InstanceBinders  map[string]InstanceBinder
	Config           brokerconfig.Config
}

// Services returns the /v2/catalog service catalog
func (kBroker *KafkaServiceBroker) Services(ctx context.Context) []brokerapi.Service {
	catalog := kBroker.loadCatalog()
	return catalog.Services
}

// Provision creates some initial Kafka topics
func (kBroker *KafkaServiceBroker) Provision(ctx context.Context, instanceID string, serviceDetails brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error) {
	spec = brokerapi.ProvisionedServiceSpec{}

	if kBroker.instanceExists(instanceID) {
		return spec, brokerapi.ErrInstanceAlreadyExists
	}

	if serviceDetails.PlanID == "" {
		return spec, errors.New("plan_id required")
	}

	planIdentifier, err := kBroker.planIdentifier(serviceDetails.PlanID)
	if err != nil {
		return spec, err
	}

	instanceCreator, ok := kBroker.InstanceCreators[planIdentifier]
	if !ok {
		return spec, errors.New("instance creator not found for plan")
	}

	err = instanceCreator.Create(instanceID)
	if err != nil {
		return spec, err
	}

	return spec, nil
}

func (kBroker *KafkaServiceBroker) planIdentifier(planID string) (string, error) {
	for _, plan := range kBroker.loadCatalog().Services[0].Plans {
		if plan.ID == planID {
			return plan.Name, nil
		}
	}
	return "", errors.New("plan_id not recognized")
}

// Deprovision deletes any topics associated with the service instance
func (kBroker *KafkaServiceBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	spec := brokerapi.DeprovisionServiceSpec{}

	for _, instanceCreator := range kBroker.InstanceCreators {
		instanceExists, _ := instanceCreator.InstanceExists(instanceID)
		if instanceExists {
			return spec, instanceCreator.Destroy(instanceID)
		}
	}
	return spec, brokerapi.ErrInstanceDoesNotExist
}

// Bind provides the information about the Kafka cluster
func (kBroker *KafkaServiceBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	binding := brokerapi.Binding{}

	for _, repo := range kBroker.InstanceBinders {
		instanceExists, _ := repo.InstanceExists(instanceID)
		if instanceExists {
			instanceCredentials, err := repo.Bind(instanceID, bindingID)
			if err != nil {
				return binding, err
			}
			credentialsMap := map[string]interface{}{
				"zkPeers":   instanceCredentials.ZookeeperPeers,
				"hostname":  instanceCredentials.KafkaHostnames,
				"topicName": instanceCredentials.TopicName,
				"uri":       fmt.Sprintf("kafka://%s/%s", instanceCredentials.KafkaHostnames, instanceCredentials.TopicName),
			}

			binding.Credentials = credentialsMap
			return binding, nil
		}
	}
	return brokerapi.Binding{}, brokerapi.ErrInstanceDoesNotExist
}

// Unbind would cancel the binding credentials
func (kBroker *KafkaServiceBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	for _, repo := range kBroker.InstanceBinders {
		instanceExists, _ := repo.InstanceExists(instanceID)
		if instanceExists {
			err := repo.Unbind(instanceID, bindingID)
			if err != nil {
				return brokerapi.ErrBindingDoesNotExist
			}
			return nil
		}
	}

	return brokerapi.ErrInstanceDoesNotExist
}

func (kBroker *KafkaServiceBroker) instanceExists(instanceID string) bool {
	for _, instanceCreator := range kBroker.InstanceCreators {
		instanceExists, _ := instanceCreator.InstanceExists(instanceID)
		if instanceExists {
			return true
		}
	}
	return false
}

// LastOperation ...
// If the broker provisions asynchronously, the Cloud Controller will poll this endpoint
// for the status of the provisioning operation.
func (kBroker *KafkaServiceBroker) LastOperation(ctx context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}

// Update would allow a service instance to be updated.
func (kBroker *KafkaServiceBroker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, nil
}
