#!/bin/bash

set -e
if [ "$PLATFORM" == "eks" ];then
 aws eks --region $AWS_DEFAULT_REGION update-kubeconfig --name $EKS_CLUSTER_NAME
fi
if [ "$MODE" == "setup" ];then
 ansible-playbook  ansible_kafka_setup.yml -vv
elif [ "$MODE" == "cleanup" ];then
 ansible-playbook ansible_kafka_cleanup.yml -vv
else
  echo "Error: Provide a valid MODE env,supported values are setup and cleanup "
  exit 1
fi
