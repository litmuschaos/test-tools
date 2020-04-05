#!/bin/sh

## mandatory arg
container_id=${CONTAINER_ID}

## optional args
memory_consumption=${MEMORY_CONSUMPTION:=500M}
duration=${DURATION:=60}
ramp_time=${RAMP_TIME:=5}

if [ ! -z "${CONTAINER_ID}" ]; then 
   
	echo "wait for the specified ramp time of ${ramp_time}s before injecting chaos"
	sleep ${ramp_time}

        echo "starting memory consumption of ${memory_consumption} Megabytes"

        docker exec ${container_id} sh -c "apt-get update && apt-get install stress-ng -y && stress-ng  --vm 1 --vm-bytes ${memory_consumption} -t ${duration}"

        echo "let chaos prevail for 30 seconds.."
        sleep 30

        echo "stopping memory chaos"
        chaos_pids=$(docker exec ${container_id} ps afx | grep "sh -c 'apt-get update'" | awk '{print$1}') 
        for i in $chaos_pids; do docker exec ${container_id} kill -9 $i; done 
else
        echo "Please provide mandatory ENV variables CONTAINER_ID & DURATION (in seconds)"
        exit 1 
fi
