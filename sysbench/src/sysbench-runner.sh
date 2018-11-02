#!/bin/bash

DB_PREFIX="sysbench"
DB_SUFFIX=`echo $(mktemp) | cut -d '.' -f 2`
DB_NAME="${DB_PREFIX}-${DB_SUFFIX}"

if [ $# -lt 2 ];
then
  echo "Usage: sh sysbench-runner.sh <db_server_ip_address> <path/to/sysbench.conf>"
  exit 1
fi 

DB_SERVER_IP=$1

T_P=($(jq -r .[] $2))
if [ $? -ne 0 ];
then
 echo -e "\nFailed to parse sysbench params, exiting\n"
 exit 1
fi 

# SYSBENCH VARS LOOKUP TABLE
#
# mysql-user     : ${T_P[0]}
# mysql-password : ${T_P[1]}
# mysql-port     : ${T_P[2]}
# db-driver      : ${T_P[3]}
# range_size     : ${T_P[4]}
# table_size     : ${T_P[5]}
# tables         : ${T_P[6]}
# threads        : ${T_P[7]}
# events         : ${T_P[8]}
# time           : ${T_P[9]}
# rand-type      : ${T_P[10]}

echo -e "\nWaiting for mysql server to start accepting connections.."
retries=30;wait_retry=30
for i in `seq 1 $retries`; do 
  mysql -h $DB_SERVER_IP -u${T_P[0]} -p${T_P[1]} -e 'status' > /dev/null 2>&1
  rc=$?
  [ $rc -eq 0 ] && break
  sleep $wait_retry 
done 

if [ $rc -ne 0 ];
then
  exit 1
  echo -e "\nFailed to connect to db server after trying for $(($retries * $wait_retry))s, exiting\n"
fi


echo -e "\nCreating database.."
mysqladmin -h $DB_SERVER_IP create $DB_NAME --user=${T_P[0]} --password=${T_P[1]} > /dev/null 2>&1
if [ $? -ne 0 ];
then
  echo -e "\nFailed to create database, exiting\n"
  exit 1
fi

echo -e "\nCreating tables.."
mysql -h $DB_SERVER_IP -u${T_P[0]} -p${T_P[1]} $DB_NAME < create_table.sql > /dev/null 2>&1
if [ $? -ne 0 ];
then
  echo -e "\nFailed to create tables, exiting\n"
  exit 1
fi

echo -e "\nLoading database.."
sysbench --mysql-host=$DB_SERVER_IP --db-driver=${T_P[3]} --mysql-user=${T_P[0]} --mysql-password=${T_P[1]} --mysql-db=$DB_NAME --range_size=${T_P[4]} --table_size=${T_P[5]} --tables=${T_P[6]} --threads=${T_P[7]} --events=${T_P[8]} --time=${T_P[9]} --rand-type=${T_P[10]} /usr/share/sysbench/oltp_read_only.lua prepare
if [ $? -ne 0 ];
then
  echo -e "\nFailed to load database, exiting\n"  
  exit 1
fi

echo -e "\nRunning benchmark.."
sysbench --mysql-host=$DB_SERVER_IP --db-driver=${T_P[3]} --mysql-user=${T_P[0]} --mysql-password=${T_P[1]} --mysql-db=$DB_NAME --range_size=${T_P[4]} --table_size=${T_P[5]} --tables=${T_P[6]} --threads=${T_P[7]} --events=${T_P[8]} --time=${T_P[9]} --rand-type=${T_P[10]} /usr/share/sysbench/oltp_read_only.lua run
if [ $? -ne 0 ];
then
  echo -e "\nFailed to run benchmark, exiting\n"
  exit 1
fi


