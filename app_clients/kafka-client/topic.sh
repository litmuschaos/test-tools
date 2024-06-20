#!/bin/bash

now="$(date +'%m/%d/%Y') $(date +"%T")"
echo "$now: [Info] Starting topic"
printf "$now: [Info] The topic information is as follow:
      Kafka Service: ${KAFKA_SERVICE}
      Kafka Port: ${KAFKA_PORT}
      Kafka Instance Name: ${KAFKA_INSTANCE_NAME}
      Topic Name: ${TOPIC_NAME}
      Replication Factor: ${REPLICATION_FACTOR}
      "
if [ ! -z "${KAFKA_INSTANCE_NAME}" ]; then
  kafka_uri=$(echo ${KAFKA_SERVICE}:${KAFKA_PORT}/${KAFKA_INSTANCE_NAME})
else
  kafka_uri=$(echo ${KAFKA_SERVICE}:${KAFKA_PORT})
fi

output=$(kafka-topics --bootstrap-server ${kafka_uri} --topic ${TOPIC_NAME} --create --partitions 1 --replication-factor ${REPLICATION_FACTOR} --if-not-exists)
echo "$now: $output"
