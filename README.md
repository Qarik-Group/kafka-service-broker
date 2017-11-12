# Kafka Service Broker

## Configuration

The following environment variables can be used to configure the broker:

* `KAFKA_HOSTNAMES` - a comma separated list of `host1:port,host2:port,host3:port`, defaults to `localhost:9092`
* `REDIS_URI` or Cloud Foundry service binding with tag `redis`, defaults to `redis://localhost:6379`

## Catalog

The default service catalog is at `data/assets/catalog.json`.

You can override the entire JSON response by setting `$BROKER_CATALOG_JSON` to the JSON string.

Alternately you can make minor adjustments to some of the globally unique attributes:

* `BROKER_SERVICE_GUID` - to change the GUID of the service
* `BROKER_SERVICE_NAME` - to change the name of the service
* `BROKER_PLAN0_GUID`, `BROKER_PLAN1_GUID` - to change the GUID of the service plan (first, second, etc)

## Development

To only clone this branch:

```
git clone https://github.com/starkandwayne/kafka-service-broker -b v2-golang-broker --single-branch
```

To run from source code:

```
go run cmd/broker/main.go
```

When adding/updating `data/assets/`, remember to run `go-bindata` to embed the changes into `data/data.go`:

```
go-bindata --pkg data -o data/data.go data/assets/...
```
