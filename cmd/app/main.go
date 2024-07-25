package main

import (
	"fmt"
	"net/http"

	"github.com/ohmyc/messaggio/internal/app/config"
	"github.com/ohmyc/messaggio/internal/app/dal"
	"github.com/ohmyc/messaggio/internal/app/server"
	"github.com/ohmyc/messaggio/pkg/processing"
	"github.com/ohmyc/messaggio/pkg/processing/request"
	"github.com/ohmyc/messaggio/pkg/processing/response"
)

func main() {
	dal_, err := dal.NewDal()
	if err != nil {
		panic(err)
	}
	if err = dal_.EnsureCreated(); err != nil {
		panic(err)
	}
	requestProducer, err := request.NewProducer(config.KafkaBootstrapServers, config.RequestTopic)
	if err != nil {
		panic(err)
	}
	responseConsumer, err := processing.NewConsumer[response.Model](config.KafkaBootstrapServers, config.ResponseTopic)
	if err != nil {
		panic(err)
	}
	go func(d *dal.Dal) {
		for {
			msg := responseConsumer.Consume()
			if err = d.UpdateProcessedMessage(msg.ID, msg.ProcessedMessage); err != nil {
				fmt.Println("Cannot update processed message:", err)
			}
		}
	}(dal_)
	s := server.NewServer(requestProducer, dal_)
	if err = http.ListenAndServe(config.AppAddress, s); err != nil {
		panic(err)
	}
}
