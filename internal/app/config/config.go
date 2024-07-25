package config

import "github.com/ohmyc/messaggio/pkg/config"

var PostgresHost = config.Env("POSTGRES_HOST")
var PostgresUser = config.Env("POSTGRES_USER")
var PostgresPassword = config.Env("POSTGRES_PASSWORD")
var KafkaBootstrapServers = config.Env("KAFKA_BOOTSTRAP_SERVERS")
var AppAddress = config.Env("APP_ADDRESS")
var RequestTopic = config.Env("KAFKA_REQUEST_TOPIC")
var ResponseTopic = config.Env("KAFKA_RESPONSE_TOPIC")
