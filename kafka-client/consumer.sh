#!/bin/bash

kafka-console-consumer --bootstrap-server ${KAFKA_SERVICE}:${KAFKA_PORT} --topic ${TOPIC_NAME} | ts '[%Y-%m-%d %H:%M:%S]' | ts 'consumer_ts:' 



