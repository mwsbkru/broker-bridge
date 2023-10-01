package config

type KafkaConfig struct {
	FromUrl           string
	FromTopic         string
	FromConsumerGroup string

	ToUrl           string
	ToTopic         string
	ToConsumerGroup string
}
