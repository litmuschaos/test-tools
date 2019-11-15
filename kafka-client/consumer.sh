#!/bin/bash

kafka-console-consumer --bootstrap-server ${KAFKA_SERVICE}:${KAFKA_PORT} --topic ${TOPIC_NAME} --timeout-ms ${KAFKA_CONSUMER_TIMEOUT} | ts '[%Y-%m-%d %H:%M:%S]' 



