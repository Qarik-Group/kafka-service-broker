package main

import (
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
	// "github.com/pivotal-cf/brokerapi/auth"

	"github.com/starkandwayne/kafka-service-broker/broker"
)

func main() {
	// brokerConfigPath := configPath()

	brokerLogger := lager.NewLogger("kafka-service-broker")
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	brokerLogger.Info("Starting Kafka service broker")

	serviceBroker := &broker.KafkaServiceBroker{
		InstanceCreators: map[string]broker.InstanceCreator{},
		InstanceBinders:  map[string]broker.InstanceBinder{},
		// Config:           config,
	}

	brokerCredentials := brokerapi.BrokerCredentials{
		Username: "broker",
		Password: "broker",
	}

	brokerAPI := brokerapi.New(serviceBroker, brokerLogger, brokerCredentials)

	port := "5000"
	http.Handle("/", brokerAPI)
	brokerLogger.Fatal("http-listen", http.ListenAndServe("0.0.0.0:"+port, nil))
}
