package app

import "github.com/ohmyc/messaggio/pkg/config"

var KafkaBootstrapServers = config.Env("KAFKA_BOOTSTRAP_SERVERS")
var AppAddress = config.Env("APP_ADDRESS")
var RequestTopic = config.Env("PROCESSING_REQUEST_TOPIC")
var ResponseTopic = config.Env("PROCESSING_RESPONSE_TOPIC")
