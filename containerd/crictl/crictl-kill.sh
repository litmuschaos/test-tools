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
retry=${Retry:-90}
delay=${Delay:-2}

###########
#  MAIN   #
###########

# Deriving the chaos iterations
C_I=$((T_C_D / CI))

#Obtain the pod ID through Pod name
pod_id=$(crictl pods | grep $A_POD | awk '{print $1}')
echo "PodID: $pod_id"

for iteration in `seq 1 $C_I`; do

  #Obtain the container ID using pod name and container name
  container_id=$(crictl ps | grep $pod_id | grep $A_CONTAINER | awk '{print $1}')
  echo "Iteration: $iteration"

  #Kill the container
  result=$(crictl stop $container_id)
  echo $result

  if [[ $result != $container_id ]]; then
    echo "Unable to kill the container $container_id"
    break
  fi

  # Obtain the container ID using pod name and container name
  for retries in `seq 1 $retry`; do
    new_container_id=$(crictl ps | grep $pod_id | grep $A_CONTAINER | awk '{print $1}')
    [ -z "$new_container_id" ] || break
    sleep $delay
  done

  if [ -z "$new_container_id" ]; then
    echo "Unable to get the new container ID"
    break
  fi

  # Check if the new container is running.
  for retries in `seq 1 $retry`; do
    status=$(crictl ps | grep $new_container_id)
    [[ "$status" == *"Running"* ]] && break
    sleep $delay
  done

  if [[ "$status" != *"Running"* ]]; then
    echo "New container is not running"
    break
  fi

  # waiting for the chaos interval
  sleep $CI
done