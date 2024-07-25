package processing

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func NewConsumer[T any](bootstrapServers, topic string) (*Consumer[T], error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	if err = consumer.Subscribe(topic, nil); err != nil {
		return nil, err
	}
	return &Consumer[T]{consumer}, nil
}

type Consumer[T any] struct {
	inner *kafka.Consumer
}

func (c *Consumer[T]) Consume() *T {
	message, _ := c.inner.ReadMessage(-1)
	r := new(T)
	err := json.Unmarshal(message.Value, r)
	if err != nil {
		panic(err)
	}
	return r
}
