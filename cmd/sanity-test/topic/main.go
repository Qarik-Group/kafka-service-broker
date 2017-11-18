package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wvanbergen/kazoo-go"
)

// Receive 'topic' plan credentials via STDIN, and perform a test upon it.

func main() {
	decoder := json.NewDecoder(os.Stdin)
	// creds := Credentials{}
	creds := make(map[string]string)
	if err := decoder.Decode(&creds); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal credentials: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Loaded credentials: %#v\n", creds)

	var topicName = creds["topicName"]
	var zkPeers = creds["zkPeers"]
	var hostname = creds["hostname"]

	errors := false
	if topicName == "" {
		fmt.Fprintf(os.Stderr, "* 'topicName' was not provided\n")
		errors = true
	}
	if zkPeers == "" {
		fmt.Fprintf(os.Stderr, "* 'zkPeers' was not provided\n")
		errors = true
	}
	if hostname == "" {
		fmt.Fprintf(os.Stderr, "* 'hostname' was not provided\n")
		errors = true
	}
	if errors {
		os.Exit(2)
	}

	zkConf := kazoo.NewConfig()
	kz, err := kazoo.NewKazooFromConnectionString(zkPeers, zkConf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "* Could not connect to Kafka: %v\n", err)
		os.Exit(3)
	}
	defer func() { _ = kz.Close() }()
	exists, err := kz.Topic(topicName).Exists()
	if err != nil {
		fmt.Fprintf(os.Stderr, "* Topic %s could not be looked up: %v\n", topicName, err)
		os.Exit(4)
	}
	if !exists {
		fmt.Fprintf(os.Stderr, "* Topic %s does not exist\n", topicName)
		os.Exit(4)
	}

}
