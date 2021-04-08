#!/bin/bash

now="$(date +'%m/%d/%Y') $(date +"%T")"
echo "$now: [Info] Starting the consumer process"
printf "$now: [Info] The consumer information is as follow:
      Kafka Service: ${KAFKA_SERVICE}
      Kafka Port: ${KAFKA_PORT}
      Topic Name: ${TOPIC_NAME}
      Kafka Consumer Timeout: ${KAFKA_CONSUMER_TIMEOUT}
      KAFKA_OPTS: ${KAFKA_OPTS}
      "
if [ ! -z "${KAFKA_OPTS}" ]; then
  output=$(kafka-console-consumer --bootstrap-server ${KAFKA_SERVICE}:${KAFKA_PORT} --topic ${TOPIC_NAME} --timeout-ms ${KAFKA_CONSUMER_TIMEOUT} --consumer.config /opt/client.properties | ts '[%Y-%m-%d %H:%M:%S]')
  echo "$now: $output"
else
  output=$(kafka-console-consumer --bootstrap-server ${KAFKA_SERVICE}:${KAFKA_PORT} --topic ${TOPIC_NAME} --timeout-ms ${KAFKA_CONSUMER_TIMEOUT} | ts '[%Y-%m-%d %H:%M:%S]')
    echo "$now: $output"
fi
echo "$now: [Info] Consumer process finished"
