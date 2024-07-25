package echo

import "github.com/ohmyc/messaggio/pkg/config"

var KafkaBootstrapServers = config.Env("KAFKA_BOOTSTRAP_SERVERS")
var RequestTopic = config.Env("KAFKA_REQUEST_TOPIC")
var ResponseTopic = config.Env("KAFKA_RESPONSE_TOPIC")
