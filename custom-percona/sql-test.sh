#!/bin/bash

DB_PREFIX="Inventory"
DB_SUFFIX=`echo $(mktemp) | cut -d '.' -f 2`
DB_NAME="${DB_PREFIX}_${DB_SUFFIX}"


echo -e "\nWaiting for mysql server to start accepting connections.."
retries=10;wait_retry=60
for i in `seq 1 $retries`; do
  mysql -uroot -pk8sDem0 -e 'status' > /dev/null 2>&1
  rc=$?
  [ $rc -eq 0 ] && break
  sleep $wait_retry
done

if [ $rc -ne 0 ];
then
  echo -e "\nFailed to connect to db server after trying for $(($retries * $wait_retry))s, exiting\n"
  exit 1
fi

## DB CREATE 
mysql -uroot -pk8sDem0 -e "CREATE DATABASE $DB_NAME;"
mysql -uroot -pk8sDem0 -e "CREATE TABLE Hardware (id INTEGER, name VARCHAR(20), owner VARCHAR(20),description VARCHAR(20));" $DB_NAME

## DB WRITE
for i in {1..100}; do 
  tvalue=`cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 18 | head -n 3`; arr=($tvalue)
  mysql -uroot -pk8sDem0 -e "INSERT INTO Hardware (id, name, owner, description) values (${i}, "${arr[0]}", "${arr[1]}", "${arr[2]}");" $DB_NAME
  sync; 
done

## DB UPDATE
for i in {1..50}; do
  mysql -uroot -pk8sDem0 -e "UPDATE Hardware SET description ="master" WHERE id = $i;" $DB_NAME 
  sync; 
done

## DB PARALLEL READ/WRITE
mysql -uroot -pk8sDem0 -e "INSERT INTO Hardware SELECT * FROM Hardware;" $DB_NAME 

## DB DESTROY
mysql -uroot -pk8sDem0 -e "DROP DATABASE $DB_NAME;"
sync; 
