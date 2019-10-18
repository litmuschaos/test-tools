# Makefile for building Containers for Storage Testing
#
#
# Reference Guide - https://www.gnu.org/software/make/manual/make.html


#
# Internal variables or constants.
# NOTE - These will be executed when any make target is invoked.
#
IS_DOCKER_INSTALLED       := $(shell which docker >> /dev/null 2>&1; echo $$?)

help:
	@echo ""
	@echo "Usage:-"
	@echo "\tmake build             -- will build openebs test containers"
	@echo "\tmake deps              -- will verify build dependencies are installed"
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
	@echo "INFO:\tverifying dependencies for OpenEBS ..."

_build_tests_vdbench_image:
	@echo "INFO: Building container image for performing vdbench tests"
	cd vdbench && docker build -t openebs/tests-vdbench .

_push_tests_vdbench_image:
	@echo "INFO: Publish container (openebs/tests-vdbench)"
	cd vdbench/buildscripts && ./push

vdbench: deps _build_tests_vdbench_image _push_tests_vdbench_image

_build_linux_utils_image:
	@echo "INFO: Building container image for linux utils"
	cd linux-utils && docker build -t openebs/linux-utils .

_push_linux_utils_image:
	@echo "INFO: Publish container (openebs/linux-utils)"
	cd linux-utils/buildscripts && ./push

linux-utils: deps _build_linux_utils_image _push_linux_utils_image

_build_tests_forkbomb_image:
	@echo "INFO: Building container image for performing forkbomb tests"
	cd forkbomb && docker build -t openebs/tests-forkbomb .

_push_tests_forkbomb_image:
	@echo "INFO: Publish container (openebs/tests-forkbomb"
	cd forkbomb/buildscripts && ./push

forkbomb: deps _build_tests_forkbomb_image _push_tests_forkbomb_image


_build_tests_fio_image:
	@echo "INFO: Building container image for performing fio tests"
	cd fio && docker build -t openebs/tests-fio .

_push_tests_fio_image:
	@echo "INFO: Publish container (openebs/tests-fio)"
	cd fio/buildscripts && ./push

fio: deps _build_tests_fio_image _push_tests_fio_image


_build_tests_chaostoolkit_image:
	@echo "INFO: Building container image for performing chaostoolkit"
	cd chaostoolkit-aws && docker build -t openebs/tests-chaostoolkit .

_push_tests_chaostoolkit_image:
	@echo "INFO: Publish container (openebs/tests-chaostoolkit)"
	cd chaostoolkit-aws/buildscripts && ./push

chaostoolkit: deps _build_tests_chaostoolkit_image _push_tests_chaostoolkit_image

_build_tests_dd_client:
	@echo "INFO: Building container image for performing dd client"
	cd dd-client && docker build -t openebs/tests-dd-client .

_push_tests_dd_client:
	@echo "INFO: Publish container (openebs/tests-dd-client)"
	cd dd-client/buildscripts && ./push

dd-client: deps _build_tests_dd_client _push_tests_dd_client

_build_tests_memleak:
	@echo "INFO: Building container image for performing dd client"
	cd memleak && docker build -t openebs/tests-memleak .

_push_tests_memleak:
	@echo "INFO: Publish container (openebs/tests-memleak)"
	cd memleak/buildscripts && ./push

memleak: deps _build_tests_memleak _push_tests_memleak

_build_tests_iometer_image:
	@echo "INFO: Building container image for performing iometer tests"
	cd iometer && docker build -t openebs/tests-iometer .

_push_tests_iometer_image:
	@echo "INFO: Publish container (openebs/tests-iometer)"
	cd iometer/buildscripts && ./push

iometer: deps _build_tests_iometer_image _push_tests_iometer_image

_build_tests_mysql_client_image:
	@echo "INFO: Building container image for performing mysql tests"
	cd mysql-client && docker build -t openebs/tests-mysql-client .

_push_tests_mysql_client_image:
	@echo "INFO: Publish container (openebs/tests-mysql-client)"
	cd mysql-client/buildscripts && ./push

mysql-client: deps _build_tests_mysql_client_image _push_tests_mysql_client_image

_build_tests_sysbench_client_image:
	@echo "INFO: Building container image for performing sysbench benchmark tests"
	cd sysbench && docker build -t openebs/sysbench-client .

_push_tests_sysbench_client_image:
	@echo "INFO: Publish container (openebs/sysbench-client)"
	cd sysbench/buildscripts && ./push

sysbench-client: deps _build_tests_sysbench_client_image _push_tests_sysbench_client_image

_build_tests_tpcc_client_image:
	@echo "INFO: Building container image for performing tpcc benchmark tests"
	cd tpcc-client && docker build -t openebs/tests-tpcc-client .

_push_tests_tpcc_client_image:
	@echo "INFO: Publish container (openebs/tests-tpcc-client)"
	cd tpcc-client/buildscripts && ./push

tpcc-client: deps _build_tests_tpcc_client_image _push_tests_tpcc_client_image

_build_tests_busybox_client_image:
	@echo "INFO: Building container image for performing busybox-liveness"
	cd busybox && docker build -t openebs/busybox-client .

_push_tests_busybox_client_image:
	@echo "INFO: Publish container (openebs/busybox-client)"
	cd busybox/buildscripts && ./push

busybox: deps _build_tests_busybox_client_image _push_tests_busybox_client_image

# busybox: deps _build_tests_busybox_client_image _push_tests_busybox_client_image

#_build_tests_custom_percona_image:
#	@echo "INFO: Building container image for integrating pmm with percona"
#	cd custom-percona && docker build -t openebs/tests-custom-percona .

#_push_tests_custom_percona_image:
#	@echo "INFO: Publish container (openebs/tests-custom-percona)"
#	cd custom-percona/buildscripts && ./push

#custom-percona: deps _build_tests_custom_percona_image _push_tests_custom_percona_image

#_build_tests_mysql_master_image:
#	@echo "INFO: Building container image for mysql-master"
#	cd mysql-master && docker build -t openebs/tests-mysql-master .

#_push_tests_mysql_master_image:
#	@echo "INFO: Publish container (openebs/tests-mysql-master)"
#	cd mysql-master/buildscripts && ./push

#mysql-master: deps _build_tests_mysql_master_image _push_tests_mysql_master_image

#_build_tests_mysql_slave_image:
#	@echo "INFO: Building container image for mysql-slave"
#	cd mysql-slave && docker build -t openebs/tests-mysql-slave .

#_push_tests_mysql_slave_image:
#	@echo "INFO: Publish container (openebs/tests-mysql-slave)"
#	cd mysql-slave/buildscripts && ./push

#mysql-slave: deps _build_tests_mysql_slave_image _push_tests_mysql_slave_image

_build_tests_mongo_client_image:
	@echo "INFO: Building container image for mongo-client"
	cd mongo-client && docker build -t openebs/tests-mongo-client .

_push_tests_mongo_client_image:
	@echo "INFO: Publish container (openebs/tests-mongo-client)"
	cd mongo-client/buildscripts && ./push

mongo-client: deps _build_tests_mongo_client_image _push_tests_mongo_client_image

_build_tests_postgres_client_image:
	@echo "INFO: Building container image for postgres-client"
	cd postgres-client && docker build -t openebs/tests-postgresql-client .

_push_tests_postgres_client_image:
	@echo "INFO: Publish container (openebs/tests-postgresql-client)"
	cd postgres-client/buildscripts && ./push

postgres-client: deps _build_tests_postgres_client_image _push_tests_postgres_client_image

_build_tests_custom_client_image:
	@echo "INFO: Building container image for custom-client"
	cd custom-client && docker build -t openebs/tests-custom-client .

_push_tests_custom_client_image:
	@echo "INFO: Publish container (openebs/tests-custom-client)"
	cd custom-client/buildscripts && ./push

custom-client: deps _build_tests_custom_client_image _push_tests_custom_client_image

_build_tests_jenkins_client_image:
	@echo "INFO: Building container image for jenkins-client"
	cd jenkins-client && docker build -t openebs/tests-jenkins-client .

_push_tests_jenkins_client_image:
	@echo "INFO: Publish container (openebs/tests-jenkins-client)"
	cd jenkins-client/buildscripts && ./push

jenkins-client: deps _build_tests_jenkins_client_image _push_tests_jenkins_client_image

_build_tests_service_liveness_image:
	@echo "INFO: Building container image for service-liveness"
	cd prometheus && docker build -t openebs/service-liveness .

_push_tests_service_liveness_image:
	@echo "INFO: Publish container (openebs/service-liveness)"
	cd prometheus/buildscripts && ./push

liveness: deps _build_tests_service_liveness_image _push_tests_service_liveness_image

_build_tests_libiscsi_image:
	@echo "INFO: Building container image for libiscsi"
	cd libiscsi && docker build -t openebs/tests-libiscsi .

_push_tests_libiscsi_image:
	@echo "INFO: Publish container (openebs/tests-libiscsi)"
	cd libiscsi/buildscripts && ./push

libiscsi: deps _build_tests_libiscsi_image _push_tests_libiscsi_image

_build_logger_image:
	@echo "INFO: Building container image for logger"
	cd logger && docker build -t openebs/logger .

_push_logger_image:
	@echo "INFO: Publish container (openebs/logger)"
	cd logger/buildscripts && ./push

logger: deps _build_logger_image _push_logger_image


_build_tests_elasticsearch_stress_image:
	@echo "INFO: Building container image for performing elasticsearch-stress tests"
	cd elasticsearch-stress && docker build -t openebs/tests-elasticsearch-stress .

_push_tests_elasticsearch_stress_image:
	@echo "INFO: Publish container (openebs/tests-elasticsearch-stress)"
	cd elasticsearch-stress/buildscripts && ./push

elasticsearch-stress: deps _build_tests_elasticsearch_stress_image _push_tests_elasticsearch_stress_image

_build_gitlab_runner_infra_image:
	@echo "INFO: Building container image for gitlab-runner-infra"
	cd gitlab-runner/buildscripts && ./build.sh
gitlab-runner: deps _build_gitlab_runner_infra_image
build: deps vdbench fio iometer mysql-client tpcc-client mongo-client jenkins-client postgres-client custom-client libiscsi logger gitlab-runner



# This is done to avoid conflict with a file of same name as the targets
# mentioned in this makefile
# Currently, help, deps build are not files in repo, but are likely candidates for addition. vdbench files exists

.PHONY: help deps build vdbench
.DEFAULT_GOAL := build
