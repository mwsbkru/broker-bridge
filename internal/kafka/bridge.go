package kafka

import (
	"context"
	"fmt"
	"github.com/mwsbkru/broker-bridge/internal/config"
	"github.com/segmentio/kafka-go"
	"log"
)

type Bridge struct {
	fromUrl           string
	fromTopic         string
	fromConsumerGroup string

	toUrl           string
	toTopic         string
	toConsumerGroup string
}

func NewBridge(cfg config.KafkaConfig) *Bridge {
	return &Bridge{
		fromUrl:           cfg.FromUrl,
		fromTopic:         cfg.FromTopic,
		fromConsumerGroup: cfg.FromConsumerGroup,
		toUrl:             cfg.ToUrl,
		toTopic:           cfg.ToTopic,
		toConsumerGroup:   cfg.ToConsumerGroup,
	}
}

func (b *Bridge) Run() error {
	// make a new reader that consumes from topic-A
	fromReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{b.fromUrl},
		GroupID:  b.fromConsumerGroup,
		Topic:    b.fromTopic,
		MaxBytes: 10e6, // 10MB
	})

	toWriter := &kafka.Writer{
		Addr:                   kafka.TCP(b.toUrl),
		Topic:                  b.toTopic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	log.Println("Start consume from: ", b.fromUrl, "/", b.fromTopic, ". Send to: ", b.toUrl, "/", b.toTopic)

	for {
		message, err := fromReader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		consumedLogMessage := fmt.Sprintf(
			"Consumed message from kafka-sender at topic/partition/offset %v/%v/%v: %s = %s",
			message.Topic,
			message.Partition,
			message.Offset,
			string(message.Key),
			string(message.Value))
		log.Println(consumedLogMessage)

		err = toWriter.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(message.Key),
				Value: []byte(message.Value),
			},
		)

		if err != nil {
			failedKey := fmt.Sprintf(
				"Failed to send message to kafka-receiver. %v. Key: %v; value: %v",
				err,
				string(message.Key),
				string(message.Value))
			log.Println(failedKey)
		} else {
			log.Println("Message sent: ", string(message.Key))
		}
	}

	if err := fromReader.Close(); err != nil {
		log.Fatal("failed to close kafka-sender:", err)
	}

	if err := toWriter.Close(); err != nil {
		log.Fatal("failed to close kafka-receiver:", err)
	}

	return nil
}
