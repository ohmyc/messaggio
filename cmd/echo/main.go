package main

import (
	"fmt"
	"strings"

	"github.com/ohmyc/messaggio/internal/echo"
	"github.com/ohmyc/messaggio/pkg/processing"
	"github.com/ohmyc/messaggio/pkg/processing/request"
	"github.com/ohmyc/messaggio/pkg/processing/response"
)

func main() {
	producer, err := response.NewProducer(echo.KafkaBootstrapServers, echo.ResponseTopic)
	if err != nil {
		panic(err)
	}
	consumer, err := processing.NewConsumer[request.Model](echo.KafkaBootstrapServers, echo.RequestTopic)
	if err != nil {
		panic(err)
	}
	for {
		m := consumer.Consume()
		if err = producer.Produce(m.ID, strings.ToUpper(m.Message)); err != nil {
			fmt.Println("Error producing message:", err)
		}
	}
}
