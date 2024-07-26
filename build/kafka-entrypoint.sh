#!/bin/bash

/etc/confluent/docker/run &

while ! kafka-topics --list --bootstrap-server localhost:9092; do
  sleep 1
done

kafka-topics --create --if-not-exists --topic processing-request  --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1
kafka-topics --create --if-not-exists --topic processing-response --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1

touch /tmp/ready.lock

wait