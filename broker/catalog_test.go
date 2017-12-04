package broker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/starkandwayne/kafka-service-broker/broker"
)

var _ = Describe("Kafka Catalog", func() {
	var kafkaBroker *broker.KafkaServiceBroker

	BeforeEach(func() {
		kafkaBroker = &broker.KafkaServiceBroker{}
	})

	Describe(".Catalog", func() {
		Context("shared kafka/zk cluster only", func() {
			It("has one service, two plans", func() {
				catalog := kafkaBroker.Catalog()
				Expect(len(catalog.Services)).To(Equal(1))
				Expect(len(catalog.Services[0].Plans)).To(Equal(2))
			})
		})
	})
})
