package broker

import (
	"encoding/json"

	"github.com/pivotal-cf/brokerapi"

	"github.com/starkandwayne/kafka-service-broker/data"
)

type Catalog struct {
	Services []brokerapi.Service
}

func (kBroker *KafkaServiceBroker) loadCatalog() (catalog Catalog) {
	catalogJSON, err := data.Asset("assets/catalog.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(catalogJSON, &catalog); err != nil {
		panic(err)
	}
	return catalog
}
