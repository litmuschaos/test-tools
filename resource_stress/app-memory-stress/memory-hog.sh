#!/bin/sh

## mandatory arg
container_id=${CONTAINER_ID}

## optional args
memory_consumption=${MEMORY_CONSUMPTION:=500}
duration=${DURATION:=60}
ramp_time=${RAMP_TIME:=10}

## Here /dev/null is a blockhole which maintains a temporary buffer in memory to write
## the chunk of data assigned in bs of the dd command. Timeout is exiting the command after
## a certain chaos duration.
chaos_cmd="timeout ${duration} dd if=/dev/zero of=/dev/null bs=${memory_consumption}M"

if [ ! -z "${CONTAINER_ID}" ]; then 

    echo "wait for the specified ramp time of ${ramp_time}s before injecting chaos"
    sleep ${ramp_time}

    echo "starting memory consumption of ${memory_consumption}Megabytes"
    echo "consuming memory..."
    docker exec ${container_id} ${chaos_cmd}
    echo "Wait for ${ramp_time} seconds for graceful termination of chaos"
    sleep ${ramp_time}

else
    echo "Please provide mandatory ENV variables CONTAINER_ID & DURATION (in seconds)"
    exit 1 
fi