#!/bin/sh

## mandatory arg
container_id=${CONTAINER_ID}

## optional args
core_count=${CORES:=1}
duration=${DURATION:=60}
ramp_time=${RAMP_TIME:=5}

## initialize chaos command
chaos_cmd="md5sum /dev/zero"

if [ ! -z "${CONTAINER_ID}" ]; then 

	echo "wait for the specified ramp time of ${ramp_time}s before injecting chaos"
	sleep ${ramp_time}

        echo "starting cpu consumption for ${core_count} cores"
        i=0
        while [ ${i} -lt ${core_count} ]; do
                echo "consuming cpu core ${i}"
                docker exec ${container_id} ${chaos_cmd} &
                i=$((i+1))
        done

        echo "let chaos prevail for ${duration} seconds.."
        sleep ${duration}

        echo "stopping cpu chaos"
        chaos_pids=$(docker exec ${container_id} ps afx | grep "${chaos_cmd}" | awk '{print$1}') 
        for i in $chaos_pids; do docker exec ${container_id} kill -9 $i; done 
else
        echo "Please provide mandatory ENV variables CONTAINER_ID & DURATION (in seconds)"
        exit 1 
fi
