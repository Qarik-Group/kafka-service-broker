package broker_test

import (
	"context"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf/brokerapi"
	"github.com/starkandwayne/kafka-service-broker/broker"
	"github.com/starkandwayne/kafka-service-broker/brokerconfig"
)

type fakeInstanceCreatorAndBinder struct {
	createErr            error
	createdInstanceIds   []string
	destroyErr           error
	destroyedInstanceIds []string
	instanceCredentials  broker.InstanceCredentials
	bindingExists        bool
}

func (fakeInstanceCreatorAndBinder *fakeInstanceCreatorAndBinder) Create(instanceID string) error {
	if fakeInstanceCreatorAndBinder.createErr != nil {
		return fakeInstanceCreatorAndBinder.createErr
	}
	fakeInstanceCreatorAndBinder.createdInstanceIds = append(fakeInstanceCreatorAndBinder.createdInstanceIds, instanceID)
	return nil
}

func (fakeInstanceCreatorAndBinder *fakeInstanceCreatorAndBinder) Destroy(instanceID string) error {
	if fakeInstanceCreatorAndBinder.destroyErr != nil {
		return fakeInstanceCreatorAndBinder.destroyErr
	}
	fakeInstanceCreatorAndBinder.destroyedInstanceIds = append(fakeInstanceCreatorAndBinder.destroyedInstanceIds, instanceID)
	return nil
}

func (fakeInstanceCreatorAndBinder *fakeInstanceCreatorAndBinder) Bind(instanceID string, bindingID string) (broker.InstanceCredentials, error) {
	return fakeInstanceCreatorAndBinder.instanceCredentials, nil
}

func (fakeInstanceCreatorAndBinder *fakeInstanceCreatorAndBinder) Unbind(instanceID string, bindingID string) error {
	if !fakeInstanceCreatorAndBinder.bindingExists {
		return errors.New("unbind error")
	}
	return nil
}

func (fakeInstanceCreatorAndBinder *fakeInstanceCreatorAndBinder) InstanceExists(instanceID string) (bool, error) {
	for _, existingInstanceID := range fakeInstanceCreatorAndBinder.createdInstanceIds {
		if instanceID == existingInstanceID {
			return true, nil
		}
	}
	return false, nil
}

var _ = Describe("Kafka SB", func() {
	ctx := context.Background()

	var kafkaBroker *broker.KafkaServiceBroker
	var someCreatorAndBinder *fakeInstanceCreatorAndBinder

	const instanceID = "instanceID"
	var topicPlanID = "4820d23c-360a-11e7-9547-d78770a33c5b"
	var planName = "topic"

	var zkPeers = "localhost:2181,localhost:2182,localhost:2183"
	var kafkaHostnames = "localhost:9092,localhost:9093,localhost:9094"

	BeforeEach(func() {
		someCreatorAndBinder = &fakeInstanceCreatorAndBinder{
			instanceCredentials: broker.InstanceCredentials{
				ZookeeperPeers: zkPeers,
				KafkaHostnames: kafkaHostnames,
				TopicName:      instanceID,
			},
		}

		kafkaBroker = &broker.KafkaServiceBroker{
			InstanceCreators: map[string]broker.InstanceCreator{
				planName: someCreatorAndBinder,
			},
			InstanceBinders: map[string]broker.InstanceBinder{
				planName: someCreatorAndBinder,
			},
			Config: brokerconfig.Config{
				KafkaConfiguration: brokerconfig.KafkaConfiguration{
					ZookeeperPeers:   zkPeers,
					ZookeeperTimeout: 1000,
					KafkaHostnames:   kafkaHostnames,
				},
			},
		}
	})

	Describe(".Provision", func() {
		Context("when the plan is recognized", func() {
			It("creates an instance", func() {
				_, err := kafkaBroker.Provision(ctx, instanceID, brokerapi.ProvisionDetails{PlanID: topicPlanID}, false)
				Expect(err).NotTo(HaveOccurred())

				Expect(len(someCreatorAndBinder.createdInstanceIds)).To(Equal(1))
				Expect(someCreatorAndBinder.createdInstanceIds[0]).To(Equal(instanceID))
			})

			Context("when the instance already exists", func() {
				BeforeEach(func() {
					_, err := kafkaBroker.Provision(ctx, instanceID, brokerapi.ProvisionDetails{PlanID: topicPlanID}, false)
					Expect(err).NotTo(HaveOccurred())
				})

				It("gives an error when trying to use the same instanceID", func() {
					_, err := kafkaBroker.Provision(ctx, instanceID, brokerapi.ProvisionDetails{PlanID: topicPlanID}, false)
					Expect(err).To(Equal(brokerapi.ErrInstanceAlreadyExists))
				})
			})

			Context("when the instance creator returns an error", func() {
				BeforeEach(func() {
					someCreatorAndBinder.createErr = errors.New("something went bad")
				})

				It("returns the same error", func() {
					_, err := kafkaBroker.Provision(ctx, instanceID, brokerapi.ProvisionDetails{PlanID: topicPlanID}, false)
					Expect(err).To(MatchError("something went bad"))
				})
			})
		})

		Context("when the plan is not recognized", func() {
			It("returns a suitable error", func() {
				_, err := kafkaBroker.Provision(ctx, instanceID, brokerapi.ProvisionDetails{PlanID: "not_a_plan_id"}, false)
				Expect(err).To(MatchError("plan_id not recognized"))
			})
		})

		Context("when the plan id is not provided", func() {
			It("returns a suitable error", func() {
				_, err := kafkaBroker.Provision(ctx, instanceID, brokerapi.ProvisionDetails{}, false)
				Expect(err).To(MatchError("plan_id required"))
			})
		})
	})

	Describe(".Deprovision", func() {
		BeforeEach(func() {
			_, err := kafkaBroker.Provision(ctx, instanceID, brokerapi.ProvisionDetails{PlanID: topicPlanID}, false)
			Expect(err).NotTo(HaveOccurred())
		})

		It("destroys the instance", func() {
			_, err := kafkaBroker.Deprovision(ctx, instanceID, brokerapi.DeprovisionDetails{}, false)
			Expect(err).NotTo(HaveOccurred())

			Expect(someCreatorAndBinder.destroyedInstanceIds).To(ContainElement(instanceID))
		})

		It("returns error if instance does not exist", func() {
			_, err := kafkaBroker.Deprovision(ctx, "non-existent", brokerapi.DeprovisionDetails{}, false)
			Expect(err).To(Equal(brokerapi.ErrInstanceDoesNotExist))
		})

		Context("when the instance creator returns an error", func() {
			BeforeEach(func() {
				someCreatorAndBinder.destroyErr = errors.New("something went bad")
			})

			It("returns the same error", func() {
				_, err := kafkaBroker.Deprovision(ctx, instanceID, brokerapi.DeprovisionDetails{}, false)
				Expect(err).To(MatchError("something went bad"))
			})
		})
	})

	Describe(".Bind", func() {
		Context("when the instance exists", func() {
			BeforeEach(func() {
				someCreatorAndBinder.Create(instanceID)
			})

			It("returns credentials", func() {
				bindingID := "bindingID"

				credentials, err := kafkaBroker.Bind(ctx, instanceID, bindingID, brokerapi.BindDetails{PlanID: topicPlanID})
				Expect(err).NotTo(HaveOccurred())

				expectedCredentials := brokerapi.Binding{
					Credentials: map[string]interface{}{
						"zkPeers":   zkPeers,
						"hostname":  kafkaHostnames,
						"topicName": instanceID,
						"uri":       fmt.Sprintf("kafka://%s/%s", kafkaHostnames, instanceID),
					},
					SyslogDrainURL:  "",
					RouteServiceURL: "",
				}

				Expect(credentials).To(Equal(expectedCredentials))
			})
		})

		Context("when the instance does not exist", func() {
			It("returns brokerapi.InstanceDoesNotExist", func() {
				bindingID := "bindingID"

				_, err := kafkaBroker.Bind(ctx, instanceID, bindingID, brokerapi.BindDetails{PlanID: topicPlanID})
				Expect(err).To(Equal(brokerapi.ErrInstanceDoesNotExist))
			})
		})
	})

	Describe(".Unbind", func() {
		BeforeEach(func() {
			someCreatorAndBinder.Create(instanceID)
			_, err := kafkaBroker.Bind(ctx, instanceID, "EXISTANT-BINDING", brokerapi.BindDetails{PlanID: topicPlanID})
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns successfully if binding existed", func() {
			someCreatorAndBinder.bindingExists = true
			err := kafkaBroker.Unbind(ctx, instanceID, "EXISTANT-BINDING", brokerapi.UnbindDetails{PlanID: topicPlanID})
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns brokerapi.ErrBindingDoesNotExist if binding did not exist", func() {
			someCreatorAndBinder.bindingExists = false
			err := kafkaBroker.Unbind(ctx, instanceID, "NON-EXISTANT-BINDING", brokerapi.UnbindDetails{PlanID: topicPlanID})
			Expect(err).To(MatchError(brokerapi.ErrBindingDoesNotExist))
		})
	})
})
