package broker

import (
	"context"
	"fmt"

	"github.com/identicalaffiliation/oms-with-events/notify-service/internal/infrastructure/config"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(cfg *config.OMSNotifyServiceConfig) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: cfg.KafkaConfig.Brokers,
			Topic:   cfg.KafkaConfig.Topic,
			GroupID: cfg.KafkaConfig.GroupID,
		}),
	}
}

func (r *KafkaConsumer) FetchMessage(ctx context.Context) (kafka.Message, error) {
	return r.reader.FetchMessage(ctx)
}

func (r *KafkaConsumer) CommitMessage(ctx context.Context, msg kafka.Message) error {
	return r.reader.CommitMessages(ctx, msg)
}

func (c *KafkaConsumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, fmt.Errorf("read message: %w", err)
	}

	return msg, nil
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
