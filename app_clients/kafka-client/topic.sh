#!/bin/bash

now="$(date +'%m/%d/%Y') $(date +"%T")"
echo "$now: [Info] Starting topic"
printf "$now: [Info] The topic information is as follow:
      Zookeeper Service: ${ZOOKEEPER_SERVICE}
      Zookeeper Port: ${ZOOKEEPER_PORT}
      Kafka Instance Name: ${KAFKA_INSTANCE_NAME}
      Topic Name: ${TOPIC_NAME}
      Replication Factor: ${REPLICATION_FACTOR}
      "
if [ ! -z "${KAFKA_INSTANCE_NAME}" ]; then
  zk_uri=$(echo ${ZOOKEEPER_SERVICE}:${ZOOKEEPER_PORT}/${KAFKA_INSTANCE_NAME})
else
  zk_uri=$(echo ${ZOOKEEPER_SERVICE}:${ZOOKEEPER_PORT})
fi

output=$(kafka-topics --zookeeper ${zk_uri} --topic ${TOPIC_NAME} --create --partitions 1 --replication-factor ${REPLICATION_FACTOR} --if-not-exists)
echo "$now: $output"
