package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/errwrap"
	"github.com/wvanbergen/kazoo-go"
)

// SanityTestSharedPlanOpts represents the 'sanity-test-topic-plan' command
type SanityTestSharedPlanOpts struct {
}

// Execute is callback from go-flags.Commander interface
func (c SanityTestSharedPlanOpts) Execute(_ []string) (err error) {
	decoder := json.NewDecoder(os.Stdin)
	// creds := Credentials{}
	creds := make(map[string]string)
	err = decoder.Decode(&creds)
	if err != nil {
		return errwrap.Wrapf("Failed to unmarshal credentials: {{err}}", err)
	}
	fmt.Printf("Loaded credentials: %#v\n", creds)

	var topicNamePrefix = creds["topicNamePrefix"]
	var zkPeers = creds["zkPeers"]
	var hostname = creds["hostname"]

	if topicNamePrefix == "" {
		return fmt.Errorf("'topicNamePrefix' was not provided")
	}
	if zkPeers == "" {
		return fmt.Errorf("'zkPeers' was not provided")
	}
	if hostname == "" {
		return fmt.Errorf("'hostname' was not provided")
	}

	zkConf := kazoo.NewConfig()
	kz, err := kazoo.NewKazooFromConnectionString(zkPeers, zkConf)
	if err != nil {
		return errwrap.Wrapf("Could not connect to Kafka: {{err}}", err)
	}
	defer func() { _ = kz.Close() }()
	exists, err := kz.Topic(topicNamePrefix).Exists()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Expected that service plan internally provisions topic %s, but could not be looked up: {{err}}", topicNamePrefix), err)
	}
	if !exists {
		return fmt.Errorf("Expected that service plan internally provisions topic %s, but does not exist", topicNamePrefix)
	}

	return
}
