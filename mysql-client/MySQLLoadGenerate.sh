#! /bin/sh

#############################################################
# This script will create an empty database and fill
# it with data continuously till the pod is alive.
#############################################################


#######################
##   FUCNTION DEF    ##       
#######################

function show_help(){
    cat << EOF
Usage : $(basename "$0") -h help
        $(basename "$0") -k <pod label key>
        $(basename "$0") -v <pod label value>
        $(basename "$0") -u <mysql username>
        $(basename "$0") -p <mysql password>
        $(basename "$0") -n <pod namespace>

-h      Display this help and exit
-k      Label key of pod, ex: app
-v      Label value of pod, ex: percona
-u      Username of MYSQL, ex: root
-p      Password of MYSQL, ex: abcd1234
-n      Namespace of Percona application pod, ex: percona

Examples: 

sh MySQLLoadGenerate.sh -k app -v percona -u root -p abcd1234 -n percona

EOF
}

# Checks whether the Percona application pod is in "Running" state
function wait_for_pod(){

    pTimeOut=$1

    c=1
    while [[ $c -lt $pTimeOut ]]; do
        # Get the pod in the specified namespace
        pName=$(kubectl get pods --no-headers -n $pod_ns -l $lk=$lv -o custom-columns=:metadata.name)

        # Check if the pod is in running state
        if [[ $(kubectl get pod $pName -n $pod_ns -o go-template --template "{{.status.phase}}") == "Running" 2>/dev/null ]]; then
            # Get pod IP
            pod_ip=$(kubectl get pod $pName -n $pod_ns -o go-template --template "{{.status.podIP}}")
            return
        else
            if [[ $c -eq 1 ]]; then
                echo -n "Waiting for application pod to be in running state"
            fi
            if [[ $(( c%5 )) -eq 0 ]]; then
                echo -n "."
            fi
        fi
        c=$(( c+1 ))
        sleep 1
    done
    echo -e "\n\e[91mUnable to bring up percona pod, exiting..\e[0m"
    exit 1
}

MySQLDump()
{
    echo -e "\n\e[96mLoadGen started!!\e[0m\n"
    while true
    do
        wait_for_pod $podTimeOut;
        mysql -u$sql_un -p$sql_pw -h $pod_ip -e "INSERT INTO Hardware select * FROM Hardware;" Inventory
        # Use a sync or sudo sync command in future
        # sync/sudo sync;
        sleep 2
    done
}

PrepareMySQL()
{
    if mysql -u$sql_un -p$sql_pw -h $pod_ip -e "USE Inventory" 2>/dev/null; then
        mysql -u$sql_un -p$sql_pw -h $pod_ip -e "DROP DATABASE Inventory;"
        echo -e "Deleted existing DataBase: \e[93mInventory\e[0m"
    fi

    mysql -u$sql_un -p$sql_pw -h $pod_ip -e "CREATE DATABASE Inventory;"
    # mysql -u$sql_un -p$sql_pw -h $pod_ip -e "CREATE DATABASE IF NOT EXISTS Inventory;"
    echo -e "Created new DataBase: \e[93mInventory\e[0m"

    mysql -u$sql_un -p$sql_pw -h $pod_ip -e \
    "CREATE TABLE Hardware (SNo VARCHAR(15),Name VARCHAR(15));" Inventory
    echo -e "Created table in DataBase: \e[93mHardware\e[0m"

    mysql -u$sql_un -p$sql_pw -h $pod_ip -e \
    "INSERT INTO Hardware (SNo,Name) VALUES ('1','John');" Inventory
    echo "Added sample values in DataBase"
}

#######################
##   TEST VARIABLES  ##
#######################

# Time to wait for Percona pod to instantiate
podTimeOut=120
# Garbage IP
pod_ip=0.0.0.0
# Garbage Key
lk=app
# Garbage Value
lv=percona
# Garbage username
sql_un=root
# Garbage password
sql_pw=abcd1234
# Garbage namespace
pod_ns=percona

#######################
## VERIFY ARGUMENTS  ##
#######################

if [[ $# -eq 0 ]]; then
    show_help
    exit 1
fi

# Obtain the input arguments
while getopts ":h:k:v:u:p:n:" option
do
    case $option in

        h)  # Display help/usage
            show_help
            exit
            ;;

        k)  # Ensure the label key is specified
            # Set the lk variable
            lk=${OPTARG}
            ;;

        v)  # Ensure the label value is specified
            # Set the lv variable
            lv=${OPTARG}
            ;;

        u)  # Ensure the MYSQL user name is specified
            # Set the sql_ns variable
            sql_un=${OPTARG}
            ;;

        p)  # Ensure the MYSQL password is specified
            # Set the sql_pw variable
            sql_pw=${OPTARG}
            ;;

        n)  # Ensure Pod namespace is specified
            # Set the pod_ns variable
            pod_ns=${OPTARG}
            ;;

        *)  # Undesired arguments
            echo "Incorrect arguments provided, please check usage"
            show_help
            exit 1
            ;;
    esac
done 

########################
##   RUN TEST STEPS   ##
########################

echo -e "\nRunning script MySQLLoadGenerator.sh\n"

# Check if Percona application pod is running
wait_for_pod $podTimeOut;

# Create empty database
PrepareMySQL;

# Start populating the database
MySQLDump; 

echo -e "\e[32mLoadGen run finished!!\e[0m"
