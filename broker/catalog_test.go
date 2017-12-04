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
		Context("override $BROKER_SERVICE_GUID", func() {
			It("has no services", func() {
				os.Setenv("BROKER_SERVICE_GUID", "XXX")
				catalog := kafkaBroker.Catalog()
				Expect(catalog.Services[0].ID).To(Equal("XXX"))
			})
		})
		Context("override $BROKER_SERVICE_NAME", func() {
			It("has no services", func() {
				os.Setenv("BROKER_SERVICE_NAME", "XXX")
				catalog := kafkaBroker.Catalog()
				Expect(catalog.Services[0].Name).To(Equal("XXX"))
			})
		})
		Context("override $BROKER_PLAN0_GUID", func() {
			It("has no services", func() {
				os.Setenv("BROKER_PLAN0_GUID", "XXX")
				catalog := kafkaBroker.Catalog()
				Expect(catalog.Services[0].Plans[0].ID).To(Equal("XXX"))
				Expect(catalog.Services[0].Plans[1].ID).ToNot(Equal("XXX"))
			})
		})
		Context("override $BROKER_PLAN1_GUID", func() {
			It("has no services", func() {
				os.Setenv("BROKER_PLAN1_GUID", "XXX")
				catalog := kafkaBroker.Catalog()
				Expect(catalog.Services[0].Plans[0].ID).ToNot(Equal("XXX"))
				Expect(catalog.Services[0].Plans[1].ID).To(Equal("XXX"))
			})
		})
	})
})
