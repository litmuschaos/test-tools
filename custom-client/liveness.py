import socket
import sys
import time
import os

# Assigning the environment variables
i_w_d = os.environ['INIT_WAIT_DELAY']   # Time period (in sec) between retries for DB init check
i_r_c = os.environ['INIT_RETRY_COUNT']      # Number of retries for DB init check
l_p_s = os.environ['LIVENESS_PERIOD_SECONDS'] # Time period (in sec) between liveness checks
l_t_s = os.environ['LIVENESS_TIMEOUT_SECONDS']  # Time period (in sec) between retries for db_connect failure
l_r_c = os.environ['LIVENESS_RETRY_COUNT'] # Number of retries after a db_connect failure before declaring liveness fail
port = os.environ['PORT']
ns = os.environ['NAMESPACE']
sv = os.environ['SERVICE']

def isOpen(ip,port):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
       s.connect((ip, int(port)))
       s.shutdown(2)
       return True
    except:
       return False


def database_check(ip,port):
    for i in range(1,int(i_r_c)):
        x=isOpen(ip,port)
        if x==False:
            isOpen(ip,port)
            print("fail", flush=True)
            # sys.stdout.flush()
        else:
            print("pass", x,flush=True)
            # sys.stdout.flush()
            break
    time.sleep(int(i_w_d))


def retryConnection(ip,port):
    print("Retrying to establish connection",flush=True)
    for y in range(1,int(l_r_c)):
       z=isOpen(ip,port)
       if z==True:
           return z
       time.sleep(int(l_t_s))
    if z==False:
        return z   


def liveness(ip,port):
    while True:    
        res=isOpen(ip,port)
        if res==True:
            print ("Liveness Running",flush=True)
            sys.stdout.flush()
        else:
            print("Liveness Failed",flush=True)
            sys.stdout.flush()
            z=retryConnection(ip,port)
            if z==False:
                print("Liveness finally failed:",flush=True)
                break
        time.sleep(int(l_p_s))

  
if __name__ == '__main__':
    url = sv+"."+ns+"."+"svc.cluster.local"
    ip= url
    port=int(port)
    database_check(ip,port)
    liveness(ip,port)
