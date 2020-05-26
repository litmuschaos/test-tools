#!/bin/sh

set -e

## mandatory arg
container_id=${CONTAINER_ID}
interval=${INTERVAL:=5}
force=${FORCE:=false}
app_ns=${APP_NS}
kill_count=${KILL_COUNT}
app_label=${APP_LABEL}
duration=${DURATION}

## Capture Current time
startTimeStamp=$(date +%s)
diffTimeStamp=0
while [ ${diffTimeStamp} -lt ${duration} ]
do
    currentTimeStamp=$(date +%s)
    diffTimeStamp=$(( $currentTimeStamp - $startTimeStamp ))
#############################################################
###############    CHECKING KILL COUNT     ##################
#############################################################

## When the kill count is not defined choose any single random pod with the given label and namesapce
if [[ -z ${kill_count} ]]; then
    echo "[Inject]: Kill a random application"
    rand_pod=$(kubectl get pod -n ${app_ns} -l ${app_label} -o=custom-columns=NAME:".metadata.name" --no-headers | shuf -n1)
    app_pod=${rand_pod}
else
    # When kill count is defined select the equal number of pod for chaos with given namespace and label
    echo "[Inject]: Starting experiment with kill count value: ${kill_count}"
    pod_list=$(kubectl get pod -n ${app_ns} -l ${app_label} -o=custom-columns=NAME:".metadata.name" --no-headers | shuf -n${kill_count})
    app_pod=${pod_list}
fi

## printing the name of application pod to be killed
echo "Name of application pod to be killed: ${app_pod}"

###########################################################
###############    FORCE POD DELETE      ##################
###########################################################

## killing the application pod forcefully if force is set to true
if [[ "${force}" == "true" ]]
then
    echo "[Inject]: Killing the application pod forcefullly"
    kubectl delete pod -n ${app_ns} --force --grace-period=0 --wait=false ${app_pod}
fi

#############################################################
###############    GRACEFUL POD DELETE     ##################
#############################################################

## killing the application pod gracefully when force is empty or force is set to false
if [[ -z "$force" || "$force" == "false" ]]
then
    echo "[Inject]: Killing the application pod gracefully"
    kubectl delete pod -n ${app_ns} ${app_pod}
fi

########################################################################
###############    CHECKING STATUS FOR RECREATION     ##################
########################################################################

echo "[Status]: Verification for the recreation of application pod"
## checking the status of pod and wait for it to come in running state
n=0
flag=0
until [ "$n" -ge 90 ]
do
    echo "[Status]: Checking the status of pods"
    pod_status=$(kubectl get pods -n ${app_ns} -l ${app_label} -o custom-columns=:.status.phase --no-headers | uniq)
    [[ ${pod_status} == "Running" ]] && break
    n=$((n+1))
    echo "pod is in ${pod_status} state please wait"
    sleep 2
    if [[ "$n" -eq 90 ]]; then
    flag=1; fi
done
if [[ $flag -eq 1 ]]; then 
echo "Application pod fails to come in running state"
exit 1; fi

## checking the status of containers and wait for it to come in running state
n=0
flag=0
until [ "$n" -ge 90 ]
do
    echo "[Status]: Checking the status of containers"
    container_status=$(kubectl get pod -n ${app_ns} -l ${app_label} --no-headers -o jsonpath='{.items[*].status.containerStatuses[*].ready}' | tr ' ' '\n' | uniq)
    [[ "${container_status}" == "true" ]] && break
    n=$((n+1)) 
    echo "pod is in ${pod_status} state please wait"
    sleep 2
    if [[ "$n" -eq 90 ]]; then
    flag=1; fi
done

if [[ $flag -eq 1 ]]; then 
    echo "Containers of application pod fails to come in running state"
    exit 1;
fi

###################################################################
###############    WAITING FOR CHAOS INTERVAL    ##################
###################################################################

## waiting for the chaos interval
echo "[Wait]: Wait for the chaos interval ${interval}s"
sleep ${interval}
## End of timestamp block
done