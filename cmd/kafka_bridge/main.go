package main

import (
	"github.com/mwsbkru/broker-bridge/internal/common"
	kafkaKonfig "github.com/mwsbkru/broker-bridge/internal/config"
	"github.com/mwsbkru/broker-bridge/internal/kafka"
)

func main() {
	var bridge common.Bridge

	kafkaCfg := kafkaKonfig.KafkaConfig{
		FromUrl:           "from_kafka:9992",
		FromConsumerGroup: "kafkaBridge",
		FromTopic:         "from_topic",
		ToUrl:             "to_kafka:9892",
		ToTopic:           "to_topic",
	}

	bridge = kafka.NewBridge(kafkaCfg)
	bridge.Run()
}
