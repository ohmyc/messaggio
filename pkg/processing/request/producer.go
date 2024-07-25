package request

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
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

func newID() string {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return id.String()
}

func newKafkaMessage(topic, id, message string) *kafka.Message {
	bytes, err := json.Marshal(Model{id, message})
	if err != nil {
		panic(err)
	}
	return &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, Value: bytes}
}

func (p *Producer) Produce(message string) (string, error) {
	id := newID()
	err := p.inner.Produce(newKafkaMessage(p.topic, id, message), nil)
	if err != nil {
		return "", err
	}
	return id, nil
}
