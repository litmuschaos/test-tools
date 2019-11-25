#!/bin/bash

if [ ! -z "${KAFKA_OPTS}" ]; then
  kafka-console-consumer --bootstrap-server ${KAFKA_SERVICE}:${KAFKA_PORT} --topic ${TOPIC_NAME} --timeout-ms ${KAFKA_CONSUMER_TIMEOUT} --consumer.config /opt/client.properties | ts '[%Y-%m-%d %H:%M:%S]' 
else
  kafka-console-consumer --bootstrap-server ${KAFKA_SERVICE}:${KAFKA_PORT} --topic ${TOPIC_NAME} --timeout-ms ${KAFKA_CONSUMER_TIMEOUT} | ts '[%Y-%m-%d %H:%M:%S]'
fi





