#!/bin/bash

i=0
while true; do 
  raw_message=$(echo "** record_index=${i}" | ts '[%Y-%m-%d %H:%M:%S]')
  echo "** producer_ts: ${raw_message}" | kafka-console-producer --broker-list ${KAFKA_SERVICE}:${KAFKA_PORT} --topic ${TOPIC_NAME}
  ((i++)) 
done
