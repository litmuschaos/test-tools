#!/bin/bash

#####################
#  VAR DEFINITION   #
#####################

: << EOF
The below variables are derived from env, with 
default values specified where the env are not present
EOF

# Delay period for MySQL server database init
I_W_D=${INIT_WAIT_DELAY:-30}

# Retry count for MySQL server database init
I_R_C=${INIT_RETRY_COUNT:-10}

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
     echo "DB_USER                 : MySQL Database User"
     echo "DB_PASSWORD             : MySQL Database Password"
     echo "DB_SVC                  : MySQL service IP/FQDN"
     echo
     echo "Set these Optional env before running the script"
     echo
     echo "INIT_WAIT_DELAY         : Wait period for MySQL server database init (default: 30s)"
     echo "INIT_RETRY_COUNT        : Retry count for MySQL server database init (default: 10)"
     echo "LIVENESS_PERIOD_SECONDS : Liveness check interval (default: 10s)"
     echo "LIVENESS_TIMEOUT_SECONDS: Liveness probe failure timeout (default: 10s)"
     echo "LIVENESS_RETRY_COUNT    : Liveness probe failure retry count (default: 3)"
     echo 
     exit  
 fi  
}

# Timestamped messages 
ts_echo()
{
 echo $1 | ts
}

# Wait (300s) until the database init completes
mysql_db_init_check()
{
 ts_echo "Waiting for mysql server to start accepting connections"

 for i in `seq 1 $I_R_C`; do
   timeout -t 5 mysql -h $1 -u$2 -p$3 -e 'status' > /dev/null 2>&1
   rc=$?
   [ $rc -eq 0 ] && break
   sleep $I_W_D
 done

 if [ $rc -ne 0 ];
 then
   ts_echo "Failed to connect to db server after trying for \
$(($I_R_C * $I_W_D))s, exiting" 
   exit 1
 fi
}

# Verify MySQL server is alive 
# Kill the query after 1s if hung/stuck (typically seen w/ disconnects) 
liveness_check()
{
 timeout -t 5 mysql -h $1 -u$2 -p$3 -e 'select 1' > /dev/null 2>&1
 rc=$?
}

###########
#  MAIN   #
###########

# Verify availability of DB credentials
if [[ -z ${DB_SVC} || -z ${DB_USER} || -z ${DB_PASSWORD} ]]; then
 usage --help;
fi

# Perform DB init check
mysql_db_init_check ${DB_SVC} ${DB_USER} ${DB_PASSWORD}

# Begin Liveness check
ts_echo "MySQL server is ready to accept connections, \
starting Liveness check" 

#Initialize the db connectivity error counter
errc=0

while true
do 
 liveness_check ${DB_SVC} ${DB_USER} ${DB_PASSWORD} 
 
 if [ $rc -ne 0 ]; then

   ts_echo "Lost connection to db server, retrying after $L_T_S sec" 

   # Increment error counter; sleep for cumulative period of 10s 
   ((errc++)); sleep $((L_T_S - 1))  

   if [ $errc -ge $L_R_C ]; then 
     ts_echo "Liveness check failed after $(($L_R_C * $L_T_S))s, exiting" 
     exit 1 
   fi

 else
   # Reset error counter if connection restored 
   ts_echo "MySQL db server is alive"; errc=0; sleep $L_P_S
 
 fi
done 


