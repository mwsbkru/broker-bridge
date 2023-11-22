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

## Test environment
For debug purposes, in folder `test-environment` prepared docker-compose.yml (based on https://github.com/conduktor/kafka-stack-docker-compose), which runs 3 independent instances of kafka.

### Prepare local kafkas
`cd test-environment && docker compose up`

### Commands for run consumer and producer for each Kafka
#### Kafka 1

consumer: `docker compose exec kafka2 kafka-console-consumer --bootstrap-server kafka1:19092 --topic topic1 --from-beginning`

producer: `docker compose exec kafka3 kafka-console-producer --bootstrap-server kafka1:19092 --topic topic1`

#### Kafka 2

consumer: `docker compose exec kafka1 kafka-console-consumer --bootstrap-server kafka2:19093 --topic topic2 --from-beginning`

producer: `docker compose exec kafka3 kafka-console-producer --bootstrap-server kafka2:19093 --topic topic2`

#### Kafka 3

consumer: `docker compose exec kafka1 kafka-console-consumer --bootstrap-server kafka3:19094 --topic topic3 --from-beginning`

producer: `docker compose exec kafka2 kafka-console-producer --bootstrap-server kafka3:19094 --topic topic3`
