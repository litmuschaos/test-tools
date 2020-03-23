#!/bin/bash

get_env () {
   tmpVar=$(echo $1)
   if [ -z "$tmpVar" ]; then
      echo "Unable to get all ENVs";
      echo "Kindly make sure to provide [BLOCK_SIZE, COUNT, NAMESPACE, MOUNT_POINT, APP_LABEL, RETRY_DURATION, RETRY_COUNT]"
      exit 1;
   fi
   echo "$tmpVar"
}

blockSize=$(get_env $BLOCK_SIZE;)    
count=$(get_env $COUNT;)
namespace=$(get_env $NAMESPACE;)
mountPoint=$(get_env $MOUNT_POINT;)
app_label=$(get_env $APP_LABEL;)
retry_duration=$(get_env $RETRY_DURATION;)
retry_count=$(get_env $RETRY_COUNT;)

#Verify that the datadir used by the templates is mounted
counter=0

# Start dd I/O 
while true
do 
   podName=$(kubectl get pod -n $namespace -l $app_label --no-headers -ojsonpath='{.items[?(@.status.phase=="Running")].metadata.name}' | head -1)
   if [ -z "$podName" ]; then
      echo "Unable to get pod Name. Retrying..."
      counter=$((counter +1))
      sleep $retry_duration
      if [ "$counter" -gt "$retry_count" ]; then
         exit 1;
      fi
   else
   echo "Writing data on $mountPoint"
   randomString=$(echo $RANDOM | tr '[0-9]' '[a-z]')
   kubectl exec $podName -n $namespace -- sh -c "cd $mountPoint && dd if=/dev/urandom of=test.$randomString bs=$blockSize count=$count"  
   sleep 10
   counter=0
   fi
done
