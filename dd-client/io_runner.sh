#!/bin/bash

#######################################################################################################################
# Script Name   : io_runner.sh         									      		
# Description   : Run dd profiles on the /datadir. 
# Creation Data : 12/10/2018                                                                                          
# Modifications : None											               		
# Script Author : Sudarshan					      
#######################################################################################################################

TEST_SIZE="512k"
TEST_COUNT="1000"

#Verify that the datadir used by the templates is mounted
if ! df -h -P | grep -q datadir > /dev/null 2>&1; then
    echo -e "datadir not mounted successfully, exiting \n"
    exit 1
fi

# Start dd I/O 
for i in {1..10}
do
   echo -e "\nRunning dd profile test with size=$TEST_SIZE ... Wait for results !!\n"
   dd if=/dev/urandom of=/datadir/f$i bs=$TEST_SIZE count=$TEST_COUNT oflag=dsync 
done  
sleep 10
rm -rf /datadir/*
