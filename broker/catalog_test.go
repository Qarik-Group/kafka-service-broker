package broker_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/starkandwayne/kafka-service-broker/broker"
)

var _ = Describe("Kafka Catalog", func() {
	var kafkaBroker *broker.KafkaServiceBroker

	BeforeEach(func() {
		os.Clearenv()
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
		Context("override via $BROKER_CATALOG_JSON", func() {
			It("has no services", func() {
				os.Setenv("BROKER_CATALOG_JSON", "{\"services\":[{\"uuid\":\"x\"}]}")
				catalog := kafkaBroker.Catalog()
				Expect(len(catalog.Services)).To(Equal(1))
				Expect(len(catalog.Services[0].Plans)).To(Equal(0))
			})
		})
	})
})
