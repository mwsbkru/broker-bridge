package config

type Config struct {
	Bridge Bridge
}

type Bridge struct {
	Kafka KafkaConfig
}

type KafkaConfig struct {
	FromUrl           string `yaml:"fromUrl"`
	FromTopic         string `yaml:"fromTopic"`
	FromConsumerGroup string `yaml:"fromConsumerGroup"`

	Receivers []KafkaReceiver `yaml:"receivers"`
}

type KafkaReceiver struct {
	ToUrl           string `yaml:"toUrl"`
	ToTopic         string `yaml:"toTopic"`
	ToConsumerGroup string `yaml:"toConsumerGroup"`
}
