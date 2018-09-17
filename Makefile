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

_build_tests_fio_image:
	@echo "INFO: Building container image for performing fio tests"
	cd fio && docker build -t openebs/tests-fio .

_push_tests_fio_image:
	@echo "INFO: Publish container (openebs/tests-fio)"
	cd fio/buildscripts && ./push

fio: deps _build_tests_fio_image _push_tests_fio_image

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

_build_tests_tpcc_client_image:
	@echo "INFO: Building container image for performing tpcc benchmark tests"
	cd tpcc-client && docker build -t openebs/tests-tpcc-client .

_push_tests_tpcc_client_image:
	@echo "INFO: Publish container (openebs/tests-tpcc-client)"
	cd tpcc-client/buildscripts && ./push

tpcc-client: deps _build_tests_tpcc_client_image _push_tests_tpcc_client_image

_build_tests_custom_percona_image:
	@echo "INFO: Building container image for integrating pmm with percona"
	cd custom-percona && docker build -t openebs/tests-custom-percona .

_push_tests_custom_percona_image:
	@echo "INFO: Publish container (openebs/tests-custom-percona)"
	cd custom-percona/buildscripts && ./push

custom-percona: deps _build_tests_custom_percona_image _push_tests_custom_percona_image

_build_tests_mysql_master_image:
	@echo "INFO: Building container image for mysql-master"
	cd mysql-master && docker build -t openebs/tests-mysql-master .

_push_tests_mysql_master_image:
	@echo "INFO: Publish container (openebs/tests-mysql-master)"
	cd mysql-master/buildscripts && ./push

mysql-master: deps _build_tests_mysql_master_image _push_tests_mysql_master_image

_build_tests_mysql_slave_image:
	@echo "INFO: Building container image for mysql-slave"
	cd mysql-slave && docker build -t openebs/tests-mysql-slave .

_push_tests_mysql_slave_image:
	@echo "INFO: Publish container (openebs/tests-mysql-slave)"
	cd mysql-slave/buildscripts && ./push

mysql-slave: deps _build_tests_mysql_slave_image _push_tests_mysql_slave_image

_build_tests_sysbench_mongo_image:
	@echo "INFO: Building container image for sysbench-mongo"
	cd sysbench-mongo && docker build -t openebs/tests-sysbench-mongo .

_push_tests_sysbench_mongo_image:
	@echo "INFO: Publish container (openebs/tests-sysbench-mongo)"
	cd sysbench-mongo/buildscripts && ./push

sysbench-mongo: deps _build_tests_sysbench_mongo_image _push_tests_sysbench_mongo_image

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

build: deps vdbench fio iometer mysql-client tpcc-client custom-percona mysql-master mysql-slave sysbench-mongo libiscsi logger 


# This is done to avoid conflict with a file of same name as the targets
# mentioned in this makefile 
# Currently, help, deps build are not files in repo, but are likely candidates for addition. vdbench files exists

.PHONY: help deps build vdbench  
.DEFAULT_GOAL := build
