#!/bin/bash

FILEPATH=${TEXTFILE_PATH:=/shared_vol}
INTERVAL=${COLLECT_INTERVAL:=10}

## calculate_pv_capacity obtains the size of a PV in bytes
function calculate_pv_capacity(){

  unit=$(echo "${size_in_spec: -1}")

  case "${unit}" in 
  
  g|gi) echo $((1024*1024*1024*$(echo $1 | cut -d "${unit}" -f 1)))
     ;;
  m|mi) echo $((1024*1024*$(echo $1 | cut -d "${unit}" -f 1)))
     ;;
  k|ki) echo $((1024*$(echo $1 | cut -d "${unit}" -f 1)))
     ;;
  b|bi) echo $1 | cut -d "${unit}" -f 1
     ;;
  *) echo 0
     ;;
  esac
}

## collect_pv_capacity_metrics collects the PV capacity metrics
function collect_pv_capacity_metrics(){
 
  ##TODO: We clear the file and then proceed to derive the metrics in the for loop below.
  ## If, the block below takes time, it may cause a few seconds of "no-metrics". 
  ## This needs to be optimized. Preferable approach is to replace values v/s recreating the file.  
  > ${FILEPATH}/pv_size.prom

  for i in ${pv_list[@]}
  do
    size_in_spec=$(kubectl get pv ${i} -o custom-columns=:spec.capacity.storage --no-headers | tr '[:upper:]' '[:lower:]')
    size_in_bytes=$(calculate_pv_capacity ${size_in_spec};)
    echo "pv_capacity_bytes{persistentvolume=\"${i}\"} ${size_in_bytes}" >> ${FILEPATH}/pv_size.prom
  done
}

## collect_pv_utilization_metrics collects the PV utilization metrics
function collect_pv_utilization_metrics(){

  ##TODO: We clear the file and then proceed to derive the metrics in the for loop below.
  ## If, the block below takes time, it may cause a few seconds of "no-metrics". 
  ## This needs to be optimized. Preferable approach is to replace values v/s recreating the file.  
  > ${FILEPATH}/pv_used.prom

  declare -a pv_mount_list=()

  for i in ${pv_list[@]}
  do
    pv_mount_list+=($(df -h | grep ${i} | awk '{print $NF}'))
  done

  echo "mount list: ${pv_mount_list[@]}"
  for i in ${pv_mount_list[@]}
  do
    ## Get mount point utilization in bytes
    mount_data=$(du -sb ${i})
    utilization=$(echo ${mount_data}| cut -d " " -f 1)
    pv_name=$(basename $(echo ${mount_data} | cut -d " " -f 2))
    echo "pv_utilization_bytes{persistentvolume=\"${pv_name}\"} ${utilization}" >> ${FILEPATH}/pv_used.prom
  done
}

while true
do
  declare -a pv_list=()

  ## Select only those PVs that are bound. Several stale PVs can exist.
  for i in $(kubectl get pv -o jsonpath='{.items[?(@.status.phase=="Bound")].metadata.name}')
  do
    pv_list+=(${i})
  done

  echo "pv_list: ${pv_list[@]}"
  collect_pv_capacity_metrics;
  collect_pv_utilization_metrics;
  sleep ${INTERVAL}
done
