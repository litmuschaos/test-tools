#!/bin/bash
#######################################################################################################################
# Script Name   : bench_runner.sh         									      		
# Description   : Run vdbench I/O using the filesystem templates on the /datadir. 
# Creation Data : 20/12/2016                                                                                          
# Modifications : None											               		
# Script Author : Karthik											      
#######################################################################################################################

TEST_TEMPLATE="file/basic-rw"
TEST_DIR="datadir"
TEST_SIZE="256m"
TEST_DURATION=0
FREE_SPACE=0
miscellaneous=0
read_only=0
write_only=0

# Function definition
function show_help(){
    cat << EOF
Usage :       $(basename "$0") --template
              $(basename "$0") --size
              $(basename "$0") --duration
              $(basename "$0") --help	
	      $(basename "$0") --available_space
	      $(basename "$0") --miscellaneous
	      $(basename "$0") --read-only
	      $(basename "$0") --write-only

-h|--help    		Display this help and exit  
--template   		Select the fio template to run 
--size	     		Provide the data sample size (in Bytes as recomonded)
--duration   		Duration (in sec)
--available_space 	Get space left at mount point
--miscellaneous		Perform writes and reads on available space on given mount point
			Ex: ./fio_runner.sh --available_space /datadir --miscellaneous /datadir(mount point)
--read-only		Perform Read on the filepath which is specified(Requires file path)
--write-only		Perform Writes using fio profile(Optional argument to perform writes)

Example: ./fio_runner.sh --template file/basic-rw --size 1024m --duration 120  

EOF
}

while [[ $# -gt 0 ]]
do 
    case $1 in
        -h|-\?|--help) 	        # Display usage summary
                       		show_help
                       		exit
                       		;;
        
        --template)    		# Optional argument to specify fio profile 
                       		if [[ -n $2 ]]; then
                       		    TEST_TEMPLATE=$2
                       		    if ! ls templates/$TEST_TEMPLATE > /dev/null 2>&1; then
                       		        echo "ERROR: Template specified does not exist"
                       		        exit 1 
                       		    fi 
                       		    shift 
                       		else
                       		    echo 'ERROR: "--template" requires a valid fio profile'
                       		    exit 1
                       		fi
                       		;;

        --size)        		# Optional argument to specify data sample size 
                       		if [[ -n $2 ]]; then
                       		    TEST_SIZE=$2
                       		    shift
                       		else
                       		    echo 'ERROR: "--size" requires a valid data sample size in MB' 
                       		    exit 1 
                       		fi 
                       		;;
         
        --duration)    		# Optional argument to specify fio run duration 
                       		if [[ -n $2 ]]; then
                       		    TEST_DURATION=$2
                       		    shift
                       		else
                       		    echo 'ERROR: "--duration" requires a valid time period in sec'
                       		    exit 1 
                       		fi 
                       		;; 
	--available_space)      # Optional argument to specify
				if [[ -n $2 ]]; then
					MOUNT_POINT=$2
					FREE_SPACE=$(df $MOUNT_POINT | awk 'FNR==2 {print $4}')
					if [ $? -ne 0 ]; then
						echo 'ERROR: Enter valid mount point'
					fi
					shift
				else
				   echo 'ERROR: "--available_space" requiers a mount directory'
				fi
				;;
	--miscellaneous)         # Perform Writes and reads based on space left
				if [[ -n $2 ]]; then
					miscellaneous=1
                                        mount_point=$2
					shift
				else
					echo 'ERROR: "--read-only" requires a file path to pass'
				fi
			       ;;
	--read-only)		# Perform Only Reads
				if [[ -n $2 ]]; then
					PATH_FILE=$2
					read_only=1
					shift
				else
					echo 'ERROR: "--read-only" requires a file path to pass'
				fi
				;;
	--write-only)		## Perform Writes (Future enhancement) 
				write_only=1
			       ;;
 
         --)           # End of all options 
                       shift 
                       break
                       ;;

         *)            # Default case: If no options, so break out of the loop
                       break
    esac
    shift

done                         


#Verify that the datadir used by the templates is mounted
if ! df -h -P | grep -q datadir > /dev/null 2>&1; then
    echo -e "datadir not mounted successfully, exiting \n"
    exit 1
fi


## Converting free space to MB
if [ $FREE_SPACE -ne 0 ]
then
	convert_to_mb=$((FREE_SPACE/20480))
fi


if [ $read_only -eq 1 ]
then
	value=$(du $PATH_FILE -b | awk '{print $1}')
	TEST_SIZE=${value}b
fi



if [ $miscellaneous -eq 1 ]; then
	if [ $convert_to_mb -ge 0 ]; then
		i=0
		while [ $i -le  `expr $convert_to_mb - 5` ]
		do
			echo "
			    [global]
			    directory=$mount_point
			    filename=missle$i

                            [test]
                            rw=write
			    bs=4k " > test_template

			echo "Write Job" 
			fio test_template --size=20m --output-format=json

			if [ $? -ne 0 ]; then
				echo "Write Failed"
			fi

		        sed 's/write/read/g' test_template

			echo "Read Job"
			fio test_template --size=20m --output-format=json

			if [ $? -ne 0 ]; then
				echo "Read Failed"
			fi
			i=$((i+1))
			
		done
		rm test_template
	else
		echo "No space Left to perform writes"
	fi
	exit 0
fi

# Start fio I/O iterating through each template file
timestamp=`date +%d%m%Y_%H%M%S`
for i in `ls templates/${TEST_TEMPLATE}`
do
   profile=$(basename $i)
   echo -e "\nRunning $profile test with size=$TEST_SIZE, runtime=$TEST_DURATION... Wait for results !!\n"
   if [ $TEST_DURATION = 0 ]
   then
        fio $i --size=$TEST_SIZE --output-format=json
   else
        fio $i --size=$TEST_SIZE --runtime=$TEST_DURATION --output-format=json
   fi 
done
