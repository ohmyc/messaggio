package response

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func NewProducer(bootstrapServers, topic string) (*Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})
	if err != nil {
		return nil, err
	}
	return &Producer{producer, topic}, nil
}

type Producer struct {
	inner *kafka.Producer
	topic string
}

func newKafkaMessage(topic, id, message string) *kafka.Message {
	bytes, err := json.Marshal(Model{id, message})
	if err != nil {
		panic(err)
	}
	return &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic}, Value: bytes}
}

func (p *Producer) Produce(id string, message string) error {
	err := p.inner.Produce(newKafkaMessage(p.topic, id, message), nil)
	if err != nil {
		return err
	}
	return nil
}
