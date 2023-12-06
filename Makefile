# Makefile for building Containers for Storage Testing
# Reference Guide - https://www.gnu.org/software/make/manual/make.html

# Internal variables or constants.
# NOTE - These will be executed when any make target is invoked.
IS_DOCKER_INSTALLED       := $(shell which docker >> /dev/null 2>&1; echo $$?)

help:
	@echo ""
	@echo "Usage:-"
	@echo "\tmake deps              -- will verify build dependencies are installed"
	@echo "\tmake <test-tool>       -- will build and push specified litmus test-tool 
	@echo ""

_build_check_docker:
	@if [ $(IS_DOCKER_INSTALLED) -eq 1 ]; \
		then echo "" \
		&& echo "ERROR:\tdocker is not installed. Please install it before build." \
		&& echo "" \
		&& exit 1; \
		fi;

deps: _build_check_docker
	@echo ""
	@echo "INFO:\tverifying dependencies for test-tools ..."

_build_tests_forkbomb_image:
	@echo "INFO: Building container image for performing forkbomb tests"
	cd resource_stress/forkbomb && docker build -t litmuschaos/forkbomb .

_push_tests_forkbomb_image:
	@echo "INFO: Publish container litmuschaos/forkbomb"
	cd resource_stress/forkbomb/buildscripts && ./push

forkbomb: deps _build_tests_forkbomb_image _push_tests_forkbomb_image

_build_tests_stress-ng_image:
	@echo "INFO: Building container image for performing stress-ng tests"
	cd resource_stress/stress-ng && docker build -t litmuschaos/stress-ng .

_push_tests_stress-ng_image:
	@echo "INFO: Publish container litmuschaos/stress-ng"
	cd resource_stress/stress-ng/buildscripts && ./push

stress-ng: deps _build_tests_stress-ng_image _push_tests_stress-ng_image

_build_tests_fio_image:
	@echo "INFO: Building container image for performing fio tests"
	cd io_tools/fio && docker build -t litmuschaos/fio .

_push_tests_fio_image:
	@echo "INFO: Publish container litmuschaos/fio"
	cd io_tools/fio/buildscripts && ./push

fio: deps _build_tests_fio_image _push_tests_fio_image

_build_tests_dd_client:
	@echo "INFO: Building container image for performing dd client"
	cd io_tools/dd && docker build -t litmuschaos/dd .

_push_tests_dd_client:
	@echo "INFO: Publish container litmuschaos/dd"
	cd io_tools/dd/buildscripts && ./push

dd: deps _build_tests_dd_client _push_tests_dd_client

_build_tests_memleak:
	@echo "INFO: Building container image for performing dd client"
	cd resource_stress/memleak && docker build -t litmuschaos/memleak .

_push_tests_memleak:
	@echo "INFO: Publish container litmuschaos/memleak"
	cd resource_stress/memleak/buildscripts && ./push

memleak: deps _build_tests_memleak _push_tests_memleak

_build_tests_mysql_client_image:
	@echo "INFO: Building container image for performing mysql tests"
	cd app_clients/mysql-client && docker build -t litmuschaos/mysql-client .

_push_tests_mysql_client_image:
	@echo "INFO: Publish container litmuschaos/mysql-client"
	cd app_clients/mysql-client/buildscripts && ./push

mysql-client: deps _build_tests_mysql_client_image _push_tests_mysql_client_image

_build_tests_sysbench_client_image:
	@echo "INFO: Building container image for performing sysbench benchmark tests"
	cd io_tools/sysbench && docker build -t litmuschaos/sysbench .

_push_tests_sysbench_client_image:
	@echo "INFO: Publish container litmuschaos/sysbench"
	cd io_tools/sysbench/buildscripts && ./push

sysbench: deps _build_tests_sysbench_client_image _push_tests_sysbench_client_image

_build_tests_tpcc_client_image:
	@echo "INFO: Building container image for performing tpcc benchmark tests"
	cd io_tools/tpcc && docker build -t litmuschaos/tpcc .

_push_tests_tpcc_client_image:
	@echo "INFO: Publish container litmuschaos/tpcc"
	cd io_tools/tpcc/buildscripts && ./push

tpcc: deps _build_tests_tpcc_client_image _push_tests_tpcc_client_image

_build_tests_mongo_client_image:
	@echo "INFO: Building container image for mongo-client"
	cd app_clients/mongo-client && docker build -t litmuschaos/mongo-client .

_push_tests_mongo_client_image:
	@echo "INFO: Publish container litmuschaos/mongo-client"
	cd app_clients/mongo-client/buildscripts && ./push

mongo-client: deps _build_tests_mongo_client_image _push_tests_mongo_client_image

_build_tests_postgres_client_image:
	@echo "INFO: Building container image for postgres-client"
	cd app_clients/postgres-client && docker build -t litmuschaos/postgresql-client .

_push_tests_postgres_client_image:
	@echo "INFO: Publish container litmuschaos/postgresql-client"
	cd app_clients/postgres-client/buildscripts && ./push

postgres-client: deps _build_tests_postgres_client_image _push_tests_postgres_client_image

_build_tests_custom_client_image:
	@echo "INFO: Building container image for custom-client"
	cd custom/custom-client && docker build -t litmuschaos/custom-client .

_push_tests_custom_client_image:
	@echo "INFO: Publish container litmuschaos/custom-client"
	cd custom/custom-client/buildscripts && ./push

custom-client: deps _build_tests_custom_client_image _push_tests_custom_client_image

_build_litmus_checker:
	@echo "INFO: Building container image for litmus-checker"
	cd custom/litmus-checker && docker build -t litmuschaos/litmus-checker .

_push_litmus_checker:
	@echo "INFO: Publish container litmuschaos/litmus-checker"
	cd custom/litmus-checker && ./buildscripts/push

litmus-checker: deps _build_litmus_checker _push_litmus_checker

_build_logger_image:
	@echo "INFO: Building container image for logger"
	cd log_utils/logger && docker build -t litmuschaos/logger .

_push_logger_image:
	@echo "INFO: Publish container litmuschaos/logger"
	cd log_utils/logger/buildscripts && ./push

logger: deps _build_logger_image _push_logger_image

_build_tests_elasticsearch_stress_image:
	@echo "INFO: Building container image for performing elasticsearch-stress tests"
	cd app_clients/elasticsearch-stress && docker build -t litmuschaos/elasticsearch-stress .

_push_tests_elasticsearch_stress_image:
	@echo "INFO: Publish container litmuschaos/elasticsearch-stress)"
	cd app_clients/elasticsearch-stress/buildscripts && ./push

elasticsearch-stress: deps _build_tests_elasticsearch_stress_image _push_tests_elasticsearch_stress_image

_build_tests_kafka_client_image:
	@echo "INFO: Building container image for kafka-liveness"
	cd app_clients/kafka-client && docker build -t litmuschaos/kafka-client .

_push_tests_kafka_client_image:
	@echo "INFO: Publish container litmuschaos/kafka-client"
	cd app_clients/kafka-client/buildscripts && ./push

kafka-client: deps _build_tests_kafka_client_image _push_tests_kafka_client_image

_build_tests_app_cpu_stress_image:
	@echo "INFO: Building container image for performing app-cpu-stress"
	cd resource_stress/app-cpu-stress && docker build -t litmuschaos/app-cpu-stress .

_push_tests_app_cpu_stress_image:
	@echo "INFO: Publish container litmuschaos/app-cpu-stress"
	cd resource_stress/app-cpu-stress/buildscripts && ./push

app-cpu-stress: deps _build_tests_app_cpu_stress_image _push_tests_app_cpu_stress_image 

_build_tests_app_memory_stress_image:
	@echo "INFO: Building container image for performing app-memory-stress"
	cd resource_stress/app-memory-stress && docker build -t litmuschaos/app-memory-stress .

_push_tests_app_memory_stress_image:
	@echo "INFO: Publish container litmuschaos/app-memory-stress"
	cd resource_stress/app-memory-stress/buildscripts && ./push

app-memory-stress: deps _build_tests_app_memory_stress_image _push_tests_app_memory_stress_image 

_build_tests_nfs_client_image:
	@echo "INFO: Building container image for performing nfs mount liveness check"
	cd app_clients/nfs-client && docker build -t litmuschaos/nfs-client .

_push_tests_nfs_client_image:
	@echo "INFO: Publish container litmuschaos/nfs-client"
	cd app_clients/nfs-client/buildscripts && ./push

nfs-client: deps _build_tests_nfs_client_image _push_tests_nfs_client_image 

_build_tests_cassandra_client_image:
	@echo "INFO: Building container image for performing cassandra liveness check"
	cd app_clients/cassandra-client && docker build -t litmuschaos/cassandra-client .

_push_tests_cassandra_client_image:
	@echo "INFO: Publish container litmuschaos/cassandra-client"
	cd app_clients/cassandra-client/buildscripts && ./push

cassandra-client: deps _build_tests_cassandra_client_image _push_tests_cassandra_client_image 

_build_tests_pod_delete_image:
	@echo "INFO: Building container image for performing pod delete chaos"
	cd experiments/pod-delete && docker build -t litmuschaos/pod-deleter .

_push_tests_pod_delete_image:
	@echo "INFO: Publish container litmuschaos/pod-deleter" 
	cd experiments/pod-delete/buildscripts && ./push

pod-delete: deps _build_tests_pod_delete_image _push_tests_pod_delete_image

_build_tests_pod_delete_go_image:
	@echo "INFO: Building container image for performing pod delete chaos"
	cd experiments/generic/pod-delete && docker build -t litmuschaos/pod-delete-helper .

_push_tests_pod_delete_go_image:
	@echo "INFO: Publish container litmuschaos/pod-delete-helper" 
	cd experiments/generic/pod-delete/buildscripts && ./push

pod-delete-go: deps _build_tests_pod_delete_go_image _push_tests_pod_delete_go_image

_build_tests_container_killer_image:

	@echo "INFO: Building container image for performing crictl container-kill"
	cd containerd/crictl && docker build -t litmuschaos/container-killer .

_push_tests_container_killer_image:
	@echo "INFO: Publish container litmuschaos/container-killer"
	cd containerd/crictl/buildscripts && ./push

container-killer: deps _build_tests_container_killer_image _push_tests_container_killer_image 

_build_tests_container_kill_go_image:
	@echo "INFO: Building container image for performing container-kill chaos"
	cd experiments/generic/container-kill && docker build -t litmuschaos/container-kill-helper .

_push_tests_container_kill_go_image:
	@echo "INFO: Publish container litmuschaos/container-kill-helper" 
	cd experiments/generic/container-kill/buildscripts && ./push

container-kill-go: deps _build_tests_container_kill_go_image _push_tests_container_kill_go_image

_build_litmus_app_deployer:
	@echo "INFO: Building container image for performing litmus-app-deployer check"
	cd custom/workflow-helper/app-deployer && docker build -t litmuschaos/litmus-app-deployer .

_push_litmus_app_deployer:
	@echo "INFO: Publish container litmuschaos/litmus-app-deployer"
	cd custom/workflow-helper/app-deployer && ./buildscripts/push

litmus-app-deployer: deps _build_litmus_app_deployer _push_litmus_app_deployer

_build_litmus_qps_cmd:
	@echo "INFO: Building container image for performing litmus-qps-cmd check"
	cd custom/workflow-helper/app-qps-test && docker build -t litmuschaos/litmus-qps-cmd .

_push_litmus_qps_cmd:
	@echo "INFO: Publish container litmuschaos/litmus-qps-cmd"
	cd custom/workflow-helper/app-qps-test && ./buildscripts/push

litmus-qps-cmd: deps _build_litmus_qps_cmd _push_litmus_qps_cmd

_build_litmus_git_app_checker:
	@echo "INFO: Building container image for performing litmus-git-app-checker check"
	cd custom/workflow-helper/app-checker && docker build -t litmuschaos/litmus-git-app-checker .

_push_litmus_git_app_checker:
	@echo "INFO: Publish container litmuschaos/litmus-git-app-checker"
	cd custom/workflow-helper/app-checker && ./buildscripts/push

litmus-git-app-checker: deps _build_litmus_git_app_checker _push_litmus_git_app_checker

_build_litmus_pg_jmeter:
	@echo "INFO: Building container image for performing litmus-pg-jmeter check"
	cd custom/workflow-helper/postgres-helper/jmeter && docker build -t litmuschaos/litmus-pg-jmeter .

_push_litmus_pg_jmeter:
	@echo "INFO: Publish container litmuschaos/litmus-pg-jmeter"
	cd custom/workflow-helper/postgres-helper/jmeter && ./buildscripts/push

litmus-pg-jmeter: deps _build_litmus_pg_jmeter _push_litmus_pg_jmeter

_build_litmus_k8s:
	@echo "INFO: Building container image for litmus-k8s"
	cd custom/k8s && docker build -t litmuschaos/k8s .

_push_litmus_k8s:
	@echo "INFO: Publish container litmuschaos/k8s"
	cd custom/k8s && ./buildscripts/push

litmus-k8s: deps _build_litmus_k8s _push_litmus_k8s

_build_litmus_curl:
	@echo "INFO: Building container image for litmus-curl"
	cd custom/curl && docker build -t litmuschaos/curl .

_push_litmus_curl:
	@echo "INFO: Publish container litmuschaos/curl"
	cd custom/curl && ./buildscripts/push

litmus-curl: deps _build_litmus_curl _push_litmus_curl

_build_litmus_argocli:
	@echo "INFO: Building container image for litmuschaos/argocli"
	cd custom/argo-server && docker build -t litmuschaos/argocli .

_push_litmus_argocli:
	@echo "INFO: Publish container litmuschaos/argocli"
	cd custom/argo-server && ./buildscripts/push

litmus-argocli: deps _build_litmus_argocli _push_litmus_argocli

_build_litmus_argo_workflow_controller:
	@echo "INFO: Building container image for litmuschaos/workflow-controller"
	cd custom/argo-workflow-controller && docker build -t litmuschaos/workflow-controller .

_push_litmus_argo_workflow_controller:
	@echo "INFO: Publish container litmuschaos/workflow-controller"
	cd custom/argo-workflow-controller && ./buildscripts/push

litmus-argo-workflow-controller: deps _build_litmus_argo_workflow_controller _push_litmus_argo_workflow_controller

_build_litmus_argo_workflow_executor:
	@echo "INFO: Building container image for litmuschaos/argoexec"
	cd custom/argo-workflow-executor && docker build -t litmuschaos/argoexec .

_push_litmus_argo_workflow_executor:
	@echo "INFO: Publish container litmuschaos/argoexec"
	cd custom/argo-workflow-executor && ./buildscripts/push

litmus-argo-workflow-executor: deps _build_litmus_argo_workflow_executor _push_litmus_argo_workflow_executor

_build_litmus_mongo:
	@echo "INFO: Building container image for litmuschaos/mongo"
	cd custom/mongo && docker build -t litmuschaos/mongo .

_push_litmus_mongo:
	@echo "INFO: Publish container litmuschaos/mongo"
	cd custom/mongo && ./buildscripts/push

litmus-mongo: deps _build_litmus_mongo _push_litmus_mongo

_build_litmus_kafka_deployer:
	@echo "INFO: Building container image for litmuschaos/kafka-deployer"
	cd custom/app-setup/kafka && docker build -t litmuschaos/kafka-deployer .

_push_litmus_kafka_deployer:
	@echo "INFO: Publish container litmuschaos/kafka-deployer"
	cd custom/app-setup/kafka/buildscripts && ./push

litmus-kafka-deployer: deps _build_litmus_kafka_deployer _push_litmus_kafka_deployer

_build_litmus_pg_load:
	@echo "INFO: Building container image for litmuschaos/litmus-pg-load"
	cd custom/workflow-helper/postgres-helper/load-test && docker build -t litmuschaos/litmus-pg-load .

_push_litmus_pg_load:
	@echo "INFO: Publish container litmuschaos/litmus-pg-load"
	cd custom/workflow-helper/postgres-helper/load-test && ./buildscripts/push

litmus-pg-load: deps _build_litmus_pg_load _push_litmus_pg_load

_build_litmus_experiment_hardened_alpine:
	@echo "INFO: Building container image for litmuschaos/experiment-alpine:latest"
	cd custom/hardened-alpine/experiment/ && docker build -t litmuschaos/experiment-alpine:latest . --build-arg TARGETARCH=amd64 --build-arg LITMUS_VERSION=1.13.8

_push_litmus_experiment_hardened_alpine:
	@echo "INFO: Publish container litmuschaos/experiment-alpine"
	cd custom/hardened-alpine/experiment/ && ./buildscripts/push

litmus-experiment-hardened-alpine: deps _build_litmus_experiment_hardened_alpine _push_litmus_experiment_hardened_alpine

_build_litmus_infra_hardened_alpine:
	@echo "INFO: Building container image for litmuschaos/infra-alpine:latest"
	cd custom/hardened-alpine/infra/ && docker build -t litmuschaos/infra-alpine:latest .

_push_litmus_infra_hardened_alpine:
	@echo "INFO: Publish container litmuschaos/infra-alpine"
	cd custom/hardened-alpine/infra/ && ./buildscripts/push

litmus-infra-hardened-alpine: deps _build_litmus_infra_hardened_alpine _push_litmus_infra_hardened_alpine

_build_litmus_mongo_utils:
	@echo "INFO: Building container image for litmuschaos/mongo-utils"
	cd custom/mongo-utils && docker build -t litmuschaos/mongo-utils .

_push_litmus_mongo_utils:
	@echo "INFO: Publish container litmuschaos/mongo-utils"
	cd custom/mongo-utils && ./buildscripts/push

litmus-mongo-utils: deps _build_litmus_mongo_utils _push_litmus_mongo_utils

_build_litmusctl:
	@echo "INFO: Building container image for litmuschaos/litmusctl"
	cd custom/litmusctl && docker build -t litmuschaos/litmusctl . --build-arg TARGETARCH=amd64 

_push_litmusctl:
	@echo "INFO: Publish container litmuschaos/litmusctl"
	cd custom/litmusctl && ./buildscripts/push

litmusctl: deps _build_litmusctl _push_litmusctl

_build_litmus_redis_load:
	@echo "INFO: Building container image for litmuschaos/litmus-redis-load"
	cd custom/workflow-helper/redis-helper/load-gen && docker build -t litmuschaos/litmus-redis-load:latest .

_push_litmus_redis_load:
	@echo "INFO: Publish container litmuschaos/litmus-kgh-loadGen"
	cd custom/workflow-helper/redis-helper/load-gen && ./buildscripts/push

litmus-redis-load: deps _build_litmus_redis_load _push_litmus_redis_load

PHONY: go-build
go-build: experiment-go-binary

experiment-go-binary:
	@echo "------------------"
	@echo "--> Build experiment go binary" 
	@echo "------------------"
	@sh build/generate_go_binary

.PHONY: docker.buildx
docker.buildx:
	@echo "------------------------------"
	@echo "--> Setting up Builder        " 
	@echo "------------------------------"
	@if ! docker buildx ls | grep -q multibuilder; then\
		docker buildx create --name multibuilder;\
		docker buildx inspect multibuilder --bootstrap;\
		docker buildx use multibuilder;\
		docker run --rm --privileged multiarch/qemu-user-static --reset -p yes;\
	fi

litmus-helm-agent: deps _build_litmus_helm_agent _push_litmus_helm_agent

_build_litmus_helm_agent:
	@echo "INFO: Building container image for litmuschaos/litmus-helm-agent"
	cd custom/litmus-helm-agent/ && docker build -t litmuschaos/litmus-helm-agent .

_push_litmus_helm_agent:
	@echo "INFO: Publish container litmuschaos/litmus-helm-agent"
	cd custom/litmus-helm-agent/ && ./buildscripts/push
