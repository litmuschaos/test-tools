import shlex
import subprocess
import sys
import time

def retstatus(name, action):
    if action == 'Backup':
        resource = "backups.velero.io"
    else:
        resource = "restore.velero.io"
    check_status = "kubectl get {} {}".format(resource,name) +  " -n velero -o=jsonpath='{.status.phase}'"
    #print(str(subprocess.check_output(backup_status,shell=True),'utf-8').strip('\n').split('\n'))
    status = str(subprocess.check_output(check_status,shell=True),'utf-8')
    return status

def runCommand(cmd, name, action):
    split_cmd = shlex.split(cmd)
    pipes = subprocess.Popen(split_cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    std_out, std_err = pipes.communicate()
    # err_msg = "%s. Code: %s" % (std_err.strip(), pipes.returncode)
    #returncode will be one non-zero if something went wrong
    if pipes.returncode != 0:
        std_err = str(std_err,'utf-8').strip().split('\n')
        _ = list(map(print, std_err))
        return
    else:
        std_out = str(std_err,'utf-8').strip().split('\n')
        _ = list(map(print, std_out))
        print("{} request submitted".format(action))
    
    #print(str(subprocess.check_output(backup_status,shell=True),'utf-8').strip('\n').split('\n'))
    backup_status = retstatus(name, action)

    if backup_status != 'InProgress' and backup_status != 'Completed':
        print("{} could not be proceeded".format(action))
        print("{} status: ".format(name), backup_status)
        return
    
    if retstatus(name, action) == 'InProgress':
        print("{} is In Progress...".format(action))

    # blah="\|/-\|/-"
    # while (retstatus(name, action) == 'InProgress'):
    #     for l in blah:
    #         sys.stdout.write(l)
    #         sys.stdout.flush()
    #         sys.stdout.write('\b')
    #         time.sleep(0.2)
    while (retstatus(name, action) == 'InProgress'):
        time.sleep(0.2)
    
    status = retstatus(name, action)

    if status == 'PartiallyFailed' or status == 'Failed':
        print("{} could not be completed!".format(action))
        print("{} status: ".format(action), status)

    if status == 'Completed':
        print("{} is done!".format(action))
        return 1
        