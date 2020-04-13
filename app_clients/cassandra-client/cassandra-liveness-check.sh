#!/bin/bash

#####################
#  VAR DEFINITION   #
#####################

: << EOF
The below variables are derived from env, with 
default values specified where the env are not present
EOF

KEYSPACE_PREFIX="Inventory"
KEYSPACE_SUFFIX=`echo $(mktemp) | cut -d '.' -f 2`
KEYSPACE_NAME="${KEYSPACE_PREFIX}_${KEYSPACE_SUFFIX}"
TABLE_PREFIX="Dataset"
TABLE_NAME="${TABLE_PREFIX}_${KEYSPACE_SUFFIX}"

# Cassandra Port
Port=${CASSANDRA_PORT:-9042}

# Replication factor for keyspace
R_F=$REPLICATION_FACTOR

# Liveness probe failure timeout 
Svc=${CASSANDRA_SVC_NAME:-cassandra}

# Liveness check interval
L_P_S=${LIVENESS_PERIOD_SECONDS:-10}

# Liveness probe failure timeout 
L_T_S=${LIVENESS_TIMEOUT_SECONDS:-10}

# Liveness probe failure retry count
L_R_C=${LIVENESS_RETRY_COUNT:-3}

###################
#   FUNCTIONS     #
###################

# Describe script usage 
usage()
{
 if [[ $1 = "--help" || $1 = "-h" ]]; then
     echo 
     echo "Usage: bash $0"
     echo
     echo "Set these MANDATORY env before running the script"
     echo 
     echo "CASSANDRA_PORT              : Cassandra Transport Port"
     echo "CASSANDRA_SVC_NAME          : Cassandra Service Name"
     echo "REPLICATION_FACTOR          : Keyspace Replication Factor"
     echo
     echo "Set these Optional env before running the script"
     echo
     echo "LIVENESS_PERIOD_SECONDS : Liveness check interval (default: 10s)"
     echo "LIVENESS_TIMEOUT_SECONDS: Liveness probe failure timeout (default: 10s)"
     echo "LIVENESS_RETRY_COUNT    : Liveness probe failure retry count (default: 3)"
     echo 
     exit  
 fi  
}

# Stores status in a action.status file
function updateStatus(){
    status=$1 
    echo "status is: ${status}"
    echo ${status} > /var/tmp/action.status
}

# Create the keyspace
create_keyspace()
{
  echo "creating the keyspace"

  for i in `seq 1 $L_R_C`; do
   cqlsh $1 $2 -e "CREATE KEYSPACE IF NOT EXISTS $KEYSPACE_NAME WITH replication = {'class' : 'SimpleStrategy', 'replication_factor' : $R_F};" > /dev/null 2>&1
   rc=$?
   [ $rc -eq 0 ] && break
   sleep $L_T_S
 done
}

# Create the table
create_table()
{
  echo "creating the table"

  for i in `seq 1 $L_R_C`; do
   cqlsh $1 $2 -e "use $KEYSPACE_NAME; CREATE TABLE $TABLE_NAME( name text PRIMARY KEY, roll int );" > /dev/null 2>&1
   rc=$?
   [ $rc -eq 0 ] && break
   sleep $L_T_S
 done
}

# Insert data into table
insert_data()
{
  echo "inserting data into table"

  for i in `seq 1 $L_R_C`; do
   cqlsh $1 $2 -e "use $KEYSPACE_NAME; INSERT INTO $TABLE_NAME (name, roll) VALUES('jamesbond',007);" > /dev/null 2>&1
   rc=$?
   [ $rc -eq 0 ] && break
   sleep $L_T_S
 done
}

# Print data from table
print_data()
{
  echo "printing data from table"

  for i in `seq 1 $L_R_C`; do
   rc=0
   cqlsh $1 $2 -e "use $KEYSPACE_NAME; SELECT * from $TABLE_NAME;" || rc=-1
   
   [ $rc -eq 0 ] && break
   sleep $L_T_S
 done
}

# Drop table
drop_table()
{
  echo "Dropping table"

  for i in `seq 1 $L_R_C`; do
   cqlsh $1 $2 -e "use $KEYSPACE_NAME; DROP TABLE $TABLE_NAME;" > /dev/null 2>&1
   rc=$?
   [ $rc -eq 0 ] && break
   sleep $L_T_S
 done
}

# Drop keyspace 
drop_keyspace()
{
  echo "Dropping keyspace"

  for i in `seq 1 $L_R_C`; do
   cqlsh $1 $2 -e "DROP KEYSPACE $KEYSPACE_NAME;" > /dev/null 2>&1
   rc=$?
   [ $rc -eq 0 ] && break
   sleep $L_T_S
 done
}

###########
#  MAIN   #
###########

# Verify availability of Required Variables
if [[ -z ${CASSANDRA_PORT} || -z ${CASSANDRA_SVC_NAME} || -z ${REPLICATION_FACTOR} ]]; then
 usage --help;
fi

# Begin Liveness check
echo "starting Liveness check" 
cycle=0

while true
do 
 updateStatus CycleInprogress;

 # creating keyspace
 create_keyspace ${Svc} ${Port}
 if [ $rc -ne 0 ]; then
   updateStatus CycleFailure;
   echo "Creation of keyspace failed after $(($L_R_C * $L_T_S))s, exiting" && exit $rc
 fi

 # creating table
 create_table ${Svc} ${Port}
 if [ $rc -ne 0 ]; then
   updateStatus CycleFailure;
   echo "Creation of table failed after $(($L_R_C * $L_T_S))s, exiting" && exit $rc
 fi

 # inserting data
 insert_data ${Svc} ${Port}
 if [ $rc -ne 0 ]; then
   updateStatus CycleFailure;
   echo "data insertion into table failed after $(($L_R_C * $L_T_S))s, exiting" && exit $rc
 fi

 # printing data
 print_data ${Svc} ${Port}
 if [ $rc -ne 0 ]; then
   updateStatus CycleFailure;
   echo "Printing of table data failed after $(($L_R_C * $L_T_S))s, exiting" && exit $rc
 fi

 # dropping table
 drop_table ${Svc} ${Port}
 if [ $rc -ne 0 ]; then
   updateStatus CycleFailure;
   echo "Dropping table failed after $(($L_R_C * $L_T_S))s, exiting" && exit $rc
 fi

 # dropping keyspace
 drop_keyspace ${Svc} ${Port}
 if [ $rc -ne 0 ]; then
   updateStatus CycleFailure;
   echo "Dropping keyspace failed after $(($L_R_C * $L_T_S))s, exiting" && exit $rc
 fi

 echo "Cycle No: $cycle is completed..."; ((cycle++))
 updateStatus CycleComplete;
 sleep $L_P_S
done 
