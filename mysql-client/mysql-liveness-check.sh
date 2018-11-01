#!/bin/bash

#####################
#  VAR DEFINITION   #
#####################

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
   mysql -h $1 -u$2 -p$3 -e 'status' > /dev/null 2>&1
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

# Usage
if [ "$#" -ne 1 ]; then
    echo 
    echo "Usage:"
    echo
    echo "$0 <db-credentials.conf>"
    echo
    echo "<db-credentials.conf> JSON file with db server,user,password"
    exit 1 
fi  

# Parse the database credentials JSON file 
db=($(jq -r .[] $1))

# DB Server Cred Lookup Table

# db_server_ip : ${db[0]} 
# db_user      : ${db[1]} 
# db_password  : ${db[2]}

# Perform DB init check
mysql_db_init_check ${db[0]} ${db[1]} ${db[2]}

# Begin Liveness check
ts_echo "MySQL server is ready to accept connections, \
starting Liveness check" 

#Initialize the db connectivity error counter
errc=0

while true
do 
 liveness_check ${db[0]} ${db[1]} ${db[2]} 
 
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


