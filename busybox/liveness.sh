#!/bin/bash

liveness_cmd='kubectl exec -it $POD_NAME -n $NAMESPACE -- sh -c "cd /busybox && dd if=/dev/urandom of=test.txt bs=4k count=1024 && echo "test" >> test.txt && cat test.txt | grep test && rm test.txt"'

retry=0

while (true); do

if [ "$retry" == "$LIVENESS_RETRY_COUNT" ]; then
  break;
fi

eval $liveness_cmd >/dev/null # '>/dev/null will supress  the output of the commnad

if [ "$?" == 0 ]; then
  echo "liveness-running"
  sleep $LIVENESS_TIMEOUT_SECONDS
else
  echo "liveness-failed"
  sleep $LIVENESS_TIMEOUT_SECONDS
  retry=$((retry+1))  
fi
done

