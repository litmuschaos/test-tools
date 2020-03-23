#!/bin/bash

i=0
while true; do
  if [ ! -z "${KAFKA_OPTS}" ]; then 	
    echo "message_index ${i} produced @ $(date -u)" | kafka-console-producer --broker-list ${KAFKA_SERVICE}:${KAFKA_PORT} --topic ${TOPIC_NAME} --producer.config /opt/client.properties
  else
    echo "message_index ${i} produced @ $(date -u)" | kafka-console-producer --broker-list ${KAFKA_SERVICE}:${KAFKA_PORT} --topic ${TOPIC_NAME}
  fi

  ((i++)) 
done
