#!/bin/bash

if [ ! -z "${KAFKA_INSTANCE_NAME}" ]; then
  zk_uri=$(echo ${ZOOKEEPER_SERVICE}:${ZOOKEEPER_PORT}/${KAFKA_INSTANCE_NAME})
else
  zk_uri=$(echo ${ZOOKEEPER_SERVICE}:${ZOOKEEPER_PORT})
fi

kafka-topics --zookeeper ${zk_uri} --topic ${TOPIC_NAME} --create --partitions 1 --replication-factor ${REPLICATION_FACTOR} --if-not-exists
