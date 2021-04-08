#!/bin/bash

now="$(date +'%m/%d/%Y') $(date +"%T")"
echo "$now: [Info] Starting the producer process"
printf "$now: [Info] The producer information is as follow:
      Kafka Service: ${KAFKA_SERVICE}
      Kafka Port: ${KAFKA_PORT}
      Topic Name: ${TOPIC_NAME}
      KAFKA_OPTS: ${KAFKA_OPTS}
      "
i=0
while true; do
  if [ ! -z "${KAFKA_OPTS}" ]; then 	
    echo "message_index ${i} produced @ $(date -u)" | kafka-console-producer --broker-list ${KAFKA_SERVICE}:${KAFKA_PORT} --topic ${TOPIC_NAME} --producer.config /opt/client.properties
  else
    echo "message_index ${i} produced @ $(date -u)" | kafka-console-producer --broker-list ${KAFKA_SERVICE}:${KAFKA_PORT} --topic ${TOPIC_NAME}
  fi

  ((i++)) 
done
echo "$now: [Info] Producer process finished"
