#!/bin/bash

now="$(date +'%m/%d/%Y') $(date +"%T")"
echo "$now: [Info] Starting topic"
printf "$now: [Info] The topic information is as follow:
      Kafka Service: ${KAFKA_SERVICE}
      Kafka Port: ${KAFKA_PORT}
      Topic Name: ${TOPIC_NAME}
      Replication Factor: ${REPLICATION_FACTOR}
      "

output=$(kafka-topics --bootstrap-server ${KAFKA_SERVICE}:${KAFKA_PORT} --topic ${TOPIC_NAME} --create --partitions 1 --replication-factor ${REPLICATION_FACTOR} --if-not-exists)
echo "$now: $output"
