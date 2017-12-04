package broker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pivotal-cf/brokerapi"

	"github.com/starkandwayne/kafka-service-broker/data"
)

// Catalog contains the service catalog returned via /v2/catalog
type Catalog struct {
	Services []brokerapi.Service
}

// Catalog imports the default /v2/catalog JSON output
// Can be overridden by environment variables
func (kBroker *KafkaServiceBroker) Catalog() (catalog *Catalog) {
	if kBroker.catalog == nil {
		catalogJSON, err := data.Asset("assets/catalog.json")
		if err != nil {
			panic(err)
		}

		catalogOverride := os.Getenv("BROKER_CATALOG_JSON")
		if catalogOverride != "" {
			if _, err := os.Stat(catalogOverride); !os.IsNotExist(err) {
				catalogJSON, err = ioutil.ReadFile(catalogOverride)
				if err != nil {
					panic(err)
				}
			} else {
				catalogJSON = []byte(catalogOverride)
			}
		}

		catalog = &Catalog{}
		if err := json.Unmarshal(catalogJSON, catalog); err != nil {
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
		kBroker.catalog = catalog
	}
	return kBroker.catalog
}
