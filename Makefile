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

_build_tests_stress_image:
	@echo "INFO: Building container image for performing stress tests"
	cd resource_stress/stress && docker build -t litmuschaos/cpu .

_push_tests_stress_image:
	@echo "INFO: Publish container litmuschaos/cpu"
	cd resource_stress/stress/buildscripts && ./push

stress: deps _build_tests_stress_image _push_tests_stress_image

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

PHONY: help deps 
.DEFAULT_GOAL := deps
