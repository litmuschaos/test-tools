import time
import subprocess
import os

# Assigning the environment variables
# Command Timeout (in sec), when showmount command got struck
command_timeout = os.environ['COMMAND_TIMEOUT']
# Time period (in sec) between retries for nfs mount check
i_w_d = os.environ['INIT_WAIT_DELAY']
# Time period (in sec) between liveness checks
l_p_s = os.environ['LIVENESS_PERIOD_SECONDS']
# Time period (in sec) between retries for nfs mount failure
l_t_s = os.environ['LIVENESS_TIMEOUT_SECONDS']
# Number of retries after a nfs mount failure before declaring liveness fail
l_r_c = os.environ['LIVENESS_RETRY_COUNT']
# Persistent volume mount using NFS
volume = os.environ['VOLUME']
# NFS provisioner service IP
nfs_svc_ip = os.environ['NFS_SVC_IP']

"""
This function checks, if the volume is present in the mounted volume list.
The volume list is generating using `showmount` command
Respose:
  True: If volume founds in the volume list
  False: Not found
"""
def isMount(ip, volume):
  commands = ["showmount", "-e", ip, "--no-headers"]                         
  volume_list=subprocess.check_output(commands, timeout=int(command_timeout)) 

  if volume_list.decode("utf-8").find(volume) >= 0:
    return True
  return False

def reCheck(ip, volume):
  """
  Retries int(l_r_c) times with int(l_t_s) secs wait between each retry
  """
  for y in range(1, int(l_r_c)):
    z = isMount(ip, volume)
    if z is True:
      return True
    if z is False:
      print("Rechecking", flush=True)
      time.sleep(int(l_t_s))
  return False

def liveness(ip, volume):
  """
  Wait int(l_p_s) secs between each retry with max retry of int(l_r_c)
  """
  while True:
    res = isMount(ip, volume)
    if res is True:
        print("volume is exported", flush=True)
    else:
        print("volume is not exported", flush=True)
        rc = reCheck(ip, volume)
        if rc is False:
          print("Liveness finally failed:", flush=True)
          break
    time.sleep(int(l_p_s))

if __name__ == '__main__':
  print('Waiting for initial wait delay ...', flush=True)
  time.sleep(int(i_w_d))
  print('Starting the liveness check ...', flush=True)
  ip = nfs_svc_ip
  liveness(ip, volume)