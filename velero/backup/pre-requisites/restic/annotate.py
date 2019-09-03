import subprocess
import os

def annotate(ns, vol, label):
    #fetch pod name
    pod_comm = "kubectl get po -n {} -l {} --no-headers".format(ns, label) + "| awk '{print $1}'"
    pod = str(subprocess.check_output(pod_comm, shell=True),'utf-8').strip('\n').split('\n')
    pod = str(subprocess.check_output(pod_comm,shell=True),'utf-8').strip('\n')
    annotate = ("kubectl -n {} annotate pod/{} backup.velero.io/backup-volumes={} --overwrite".format(ns,pod,vol)).strip('/n')
    os.system(annotate)
    