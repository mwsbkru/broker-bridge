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

	receivers []Receiver
}

type Receiver struct {
	toUrl           string
	toTopic         string
	toConsumerGroup string
}

func NewBridge(cfg config.KafkaConfig) *Bridge {
	receivers := make([]Receiver, len(cfg.Receivers))

	for i, receiverCfg := range cfg.Receivers {
		receivers[i] = Receiver{
			toUrl:           receiverCfg.ToUrl,
			toTopic:         receiverCfg.ToTopic,
			toConsumerGroup: receiverCfg.ToConsumerGroup,
		}
	}
	return &Bridge{
		fromUrl:           cfg.FromUrl,
		fromTopic:         cfg.FromTopic,
		fromConsumerGroup: cfg.FromConsumerGroup,
		receivers:         receivers,
	}
}

func (b *Bridge) Run(ctx context.Context) error {
	// make a new reader that consumes from topic-A
	fromReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{b.fromUrl},
		GroupID:  b.fromConsumerGroup,
		Topic:    b.fromTopic,
		MaxBytes: 10e6, // 10MB
	})

	writers := make([]*kafka.Writer, len(b.receivers))

	for i, receiver := range b.receivers {
		writers[i] = &kafka.Writer{
			Addr:                   kafka.TCP(receiver.toUrl),
			Topic:                  receiver.toTopic,
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		}
	}

	log.Printf("Start consume from: %v/%v (consumer group: %v).\n", b.fromUrl, b.fromTopic, b.fromConsumerGroup)

	for {
		message, err := fromReader.ReadMessage(ctx)
		var processedMessage string
		if len(message.Key) > 0 {
			processedMessage = string(message.Key)
		} else {
			processedMessage = "key not set"
		}

		if err != nil {
			connectToFromKafkaErrText := fmt.Sprintf("Reading from \"From Kafka\": %v", err)
			log.Println(connectToFromKafkaErrText)
			break
		}
		consumedLogMessage := fmt.Sprintf(
			"Consumed message from kafka-sender at topic/partition/offset %v/%v/%v. Key: %s. Message: %s.",
			message.Topic,
			message.Partition,
			message.Offset,
			processedMessage,
			string(message.Value))
		log.Println(consumedLogMessage)

		for i, writer := range writers {
			err = writer.WriteMessages(ctx,
				kafka.Message{
					Key:   []byte(processedMessage),
					Value: message.Value,
				},
			)

			if err != nil {
				failedKey := fmt.Sprintf(
					"Writing to \"To Kafka\". %v. Key: %v. value: %v. Receiver url: %v. Topic: %v",
					err,
					string(message.Key),
					string(message.Value),
					b.receivers[i].toUrl,
					b.receivers[i].toTopic)
				log.Println(failedKey)
			} else {
				messegeSendedText := fmt.Sprintf(
					"Message sent: %v. Receiver url: %v. Topic: %v",
					processedMessage,
					b.receivers[i].toUrl,
					b.receivers[i].toTopic)
				log.Println(messegeSendedText)
			}
		}

		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}

	if err := fromReader.Close(); err != nil {
		log.Fatal("failed to close kafka-sender:", err)
	}

	for _, writer := range writers {
		if err := writer.Close(); err != nil {
			log.Println("failed to close kafka-receiver:", err)
		}
	}

	return nil
}
