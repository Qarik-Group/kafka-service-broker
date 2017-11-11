# Kafka Service Broker

To run from source code:

```
go run cmd/broker/main.go
```

When adding/updating `data/assets/`, remember to run `go-bindata` to embed the changes into `data/data.go`:

```
go-bindata --pkg data -o data/data.go data/assets/...
```
