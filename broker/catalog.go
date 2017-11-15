package broker

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pivotal-cf/brokerapi"

	"github.com/starkandwayne/kafka-service-broker/data"
)

// Catalog contains the service catalog returned via /v2/catalog
type Catalog struct {
	Services []brokerapi.Service
}

func (kBroker *KafkaServiceBroker) loadCatalog() (catalog Catalog) {
	catalogJSON, err := data.Asset("assets/catalog.json")
	if err != nil {
		panic(err)
	}
	if os.Getenv("BROKER_CATALOG_JSON") != "" {
		catalogJSON = []byte(os.Getenv("BROKER_CATALOG_JSON"))
	}

	if err := json.Unmarshal(catalogJSON, &catalog); err != nil {
		panic(err)
	}
	if os.Getenv("BROKER_SERVICE_GUID") != "" {
		catalog.Services[0].ID = os.Getenv("BROKER_SERVICE_GUID")
	}
	if os.Getenv("BROKER_SERVICE_NAME") != "" {
		catalog.Services[0].Name = os.Getenv("BROKER_SERVICE_NAME")
	}
	for i := range catalog.Services[0].Plans {
		if os.Getenv(fmt.Sprintf("BROKER_PLAN%d_GUID", i)) != "" {
			catalog.Services[0].Plans[i].ID = os.Getenv(fmt.Sprintf("BROKER_PLAN%d_GUID", i))
		}
	}

	if os.Getenv("CATALOG_DOCUMENTATION_URL") != "" {
		catalog.Services[0].Metadata.DocumentationUrl = os.Getenv("CATALOG_DOCUMENTATION_URL")
	}

	return catalog
}
