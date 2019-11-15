#!/bin/bash

i=0
while true; do 
  echo "message_index ${i} produced @ $(date -u)" | kafka-console-producer --broker-list ${KAFKA_SERVICE}:${KAFKA_PORT} --topic ${TOPIC_NAME}
  ((i++)) 
done
