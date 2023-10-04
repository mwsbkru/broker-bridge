# Broker-bridge
Tool for forward messages from one kafka server to another. May be helpful for developing products, running in different Docker environments with his own kafka container.

## Run from sources
Make sure both kafka servers runs and accessible from host, running bridge.
Edit `cmd/kafka_bridge/config.yaml` according source and receiver kafka servers.
Run the following command.

````
cd cmd/kafka_bridge && go run main.go
````

## Run from binary
Make sure both kafka servers runs and accessible from host, running bridge.
Build executable file (`cd cmd/kafka_bridge && go build`).
Copy `cmd/kafka_bridge/config.yaml` to directory with `kafka_bridge` executable and edit it according source and receiver kafka servers.
Run the following command in directory with executable and config.

````
./kafka_bridge
````


