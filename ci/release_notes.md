# Rewritten in Golang

v2 of starkandwayne/kafka-service-broker is a rewrite in Golang. You can find the previous Spring version at https://github.com/starkandwayne/kafka-service-broker/tree/v1-spring-broker. The `master` git branch has no shared history with `v1-spring-broker` branch, but we'll retain the two versions in the same repository for now.

The functionality of this broker is similar to v1. It operates against a singular, shared ZooKeeper/Kafka cluster. It provides a service `starkandwayne-kafka` and single plan `topic`.

The service plan `topic` offers the same `hostname` (comma-separated list of kafka brokers) and `topicName` credentials as v1. In addition, `credentials` now includes `zkPeers` - the comma-separated list of zookeeper peers. Some Kafka client libraries prefer to dynamically lookup the Kafka brokers from ZooKeeper and this is now possible.

The README has been updated to discuss configuration via environment variables.

* Use `$BROKER_USERNAME` and `$BROKER_PASSWORD` to set the basic auth credentials.
* Provide `$ZOOKEEPER_PEERS` to allow the broker to discover the Kafka brokers/nodes.
* The v2 broker now supports customising the entire `/v2/catalog` via the `$BROKER_CATALOG_JSON` environment variable, or you can modify the GUIDs/service name using other variables.
