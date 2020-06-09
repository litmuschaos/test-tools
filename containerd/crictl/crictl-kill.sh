#!/bin/bash

#####################
#  VAR DEFINITION   #
#####################

: << EOF
The below variables are derived from env, with 
default values specified where the env are not present
EOF

A_CONTAINER=$APP_CONTAINER
A_POD=$APP_POD
CI=$CHAOS_INTERVAL
T_C_D=$TOTAL_CHAOS_DURATION
retry=${RETRY:-90}
delay=${DELAY:-2}
C_NS=$CHAOS_NAMESPACE
Engine=$CHAOS_ENGINE
E_UID=$CHAOS_UID
C_POD=$POD_NAME
C_I=$ITERATIONS

###########
#  MAIN   #
###########

#Obtain the pod ID through Pod name
pod_id=$(crictl pods | grep $A_POD | awk '{print $1}')
echo "App Pod ID: $pod_id"

startTimeStamp=$(date +%s)

for iteration in `seq 1 $C_I`; do

  #Obtain the container ID using pod name and container name
  container_id=$(crictl ps | grep $pod_id | grep $A_CONTAINER | awk '{print $1}')
  echo "Iteration: $iteration"

  if [[ ! -z ${Engine} ]]; then
    # get the current ts
    TS=$( date -u +"%Y-%m-%dT%H:%M:%SZ")
    echo "Timestamp: $TS"
    #Creating ChaosInject Event
    jinja2 -D engine_ns=$C_NS -D chaos_pod=$C_POD -D app_pod=$A_POD -D app_container=$A_CONTAINER -D ts=$TS -D engine_name=$Engine -D val=$iteration -D engine_uid=$E_UID event.yaml > rendered_event.yaml
  
    #apply the events
    kubectl apply -f rendered_event.yaml
  fi

  #Kill the container
  result=$(crictl stop $container_id)
  echo "App Container ID $result"

  if [[ $result != $container_id ]]; then
    echo "Unable to kill the application container $container_id"
    exit 1
  fi

  # Obtain the container ID using pod name and container name
  for retries in `seq 1 $retry`; do
    new_container_id=$(crictl ps | grep $pod_id | grep $A_CONTAINER | awk '{print $1}')
    [ -z "$new_container_id" ] || break
    sleep $delay
  done

  if [ -z "$new_container_id" ]; then
    echo "Unable to get the new application container ID"
    exit 1
  fi

  # Check if the new container is running.
  for retries in `seq 1 $retry`; do
    status=$(crictl ps | grep $new_container_id)
    [[ "$status" == *"Running"* ]] && break
    sleep $delay
  done

  if [[ "$status" != *"Running"* ]]; then
    echo "Restarted container is not in running state"
    exit 1
  fi

  # waiting for the chaos interval
  sleep $CI

  currentTimeStamp=$(date +%s)

  diffTimeStamp="$(($currentTimeStamp-$startTimeStamp))"
  
  if [[ $diffTimeStamp -ge $T_C_D ]]; then
    echo "terminating the execution after $diffTimeStamp s"
    exit 0
  fi
done