#!/bin/bash

liveness_cmd='curl $SERVICE_NAME.$NAMESPACE.svc.cluster.local:9200'

eval $liveness_cmd

retry=0

while (true); do

if [ "$retry" == "$LIVENESS_RETRY_COUNT" ]; then
  break;
fi

eval $liveness_cmd >/dev/null # '>/dev/null will supress the output of the command

if [ "$?" == 0 ]; then
  echo "liveness-running"
  sleep $LIVENESS_TIMEOUT_SECONDS
else
  echo "livenes-failed"
  sleep $LIVENESS_TIMEOUT_SECONDS
  retry=$((retry+1))
fi
done

