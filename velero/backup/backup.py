#!/usr/bin/env python3
""" This script executes velero backup """
import argparse
import sys
import re
import datetime
import argparse
import subprocess
import os
import time
sys.path.append('../newCommand/backup/')
import newBackupCommand 
sys.path.append('../runCommand/')
import runCommand
sys.path.append('pre-requisites/restic')
import annotate

#python3 backup1.py -n testp -ns=minio,minio-lpd --selector='!openebs.io/controller,!openebs.io/replica' --storage-location=gcp --volume-snapshot-locations=gcp
def time_diff(date2, date1, format):
    ''' Calculates time difference by backup '''
    diff = datetime.datetime.strptime(date2, format) - datetime.datetime.strptime(date1, format)
    return diff
	
def time_in_mins(backup_name):
    ''' Calculates time taken by backup to complete    '''
    datetimeFormat = '%Y-%m-%dT%H:%M:%SZ'
    completiontime = "kubectl get backups {} -n velero".format(backup_name) + " -o=jsonpath='{.status.completionTimestamp}'"
    starttime = "kubectl get backups {} -n velero".format(backup_name) +  " -o=jsonpath='{.status.startTimestamp}'"
    completiontime = str(subprocess.check_output(completiontime, shell=True),'utf-8')
    starttime = str(subprocess.check_output(starttime, shell=True),'utf-8')
    diff = time_diff(completiontime, starttime, datetimeFormat)
    sec = diff.seconds 
    if sec == 60:
        print("Backup completion time: 1min")
    elif sec > 60:
        min = sec//60
        sec = sec % 60
        print("Backup completion time: {} min {} secs".format(min,sec))
    else:
        print("Backup completion time: {} secs".format(sec))
   
def main():
    parser = argparse.ArgumentParser()

    parser.add_argument(
        "-n", "--name",
        help="name")
    
    parser.add_argument(
        "-sch", "--schedule",
        help="schedule time interval")

    parser.add_argument(
        "-ns", "--include-namespaces",
        nargs='+', 
        help="provide comma seperated list of namespaces you want to backup")

    parser.add_argument(
        "-loc", "--volume-snapshot-locations", 
        help="provide volume snapshot location")
    
    parser.add_argument("--volns", 
        nargs='+', 
        help="provide comma seperated namespace:vol pair")

    parser.add_argument("--selector",
        help="label selector")

    parser.add_argument("--storage-location",
        help="backup location")

    parser.add_argument("-ttl", 
        "--default-backup-ttl", 
        help="The amount of time before this backup is eligible for garbage collection. If not specified, a default value of 30 days will be used. The default can be configured on the velero server by passing the flag ttl: 24h0m0s")
    
    parser.add_argument("--app",
        help="app for which backup is performed")

    # <class 'argparse.Namespace'>
    args = parser.parse_args()
    name = args.name
    command = newBackupCommand.newBackupCommand(args)  
    print("Velero commnad: ", command)
    action = "Backup"
    
    if args.volns != None:
        input = args.volns[0].split(',')
        # label will later be taken by litmus testvars
        label = "name=" + args.app 
        for i in input:
            ns = i.split(':')[0]
            vol = i.split(':')[1]
            annotate.annotate(ns, vol, label)
        print("Pre-requisites for restic met successfully")

    status = runCommand.runCommand(command, name, action)
    if status == 1:
        time_in_mins(name)
    
if __name__ == '__main__':
    main()