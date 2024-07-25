package main

import (
	"fmt"
	"net/http"

	"github.com/ohmyc/messaggio/internal/app"
	"github.com/ohmyc/messaggio/pkg/processing"
	"github.com/ohmyc/messaggio/pkg/processing/request"
	"github.com/ohmyc/messaggio/pkg/processing/response"
)

func main() {
	dal, err := app.NewDal()
	if err != nil {
		panic(err)
	}
	if err = dal.EnsureCreated(); err != nil {
		panic(err)
	}
	requestProducer, err := request.NewProducer(app.KafkaBootstrapServers, app.RequestTopic)
	if err != nil {
		panic(err)
	}
	responseConsumer, err := processing.NewConsumer[response.Model](app.KafkaBootstrapServers, app.ResponseTopic)
	if err != nil {
		panic(err)
	}
	go func(d *app.Dal) {
		for {
			msg := responseConsumer.Consume()
			if err = d.UpdateProcessedMessage(msg.ID, msg.ProcessedMessage); err != nil {
				fmt.Println("Cannot update processed message:", err)
			}
		}
	}(dal)
	s := app.NewServer(requestProducer, dal)
	if err = http.ListenAndServe(app.AppAddress, s); err != nil {
		panic(err)
	}
}
