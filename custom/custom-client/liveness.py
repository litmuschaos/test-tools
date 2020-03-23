import socket
import sys
import time
import os

# Assigning the environment variables
# Time period (in sec) between retries for DB init check
i_w_d = os.environ['INIT_WAIT_DELAY']
# Number of retries for DB init check
i_r_c = os.environ['INIT_RETRY_COUNT']
# Time period (in sec) between liveness checks
l_p_s = os.environ['LIVENESS_PERIOD_SECONDS']
# Time period (in sec) between retries for db_connect failure
l_t_s = os.environ['LIVENESS_TIMEOUT_SECONDS']
# Number of retries after a db_connect failure before declaring liveness fail
l_r_c = os.environ['LIVENESS_RETRY_COUNT']
port = os.environ['PORT']
ns = os.environ['NAMESPACE']
sv = os.environ['SERVICE']


def isOpen(ip, port):
    """
    Socket connection to url:port
    """
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        s.connect((ip, int(port)))
        s.shutdown(2)  # Disables the socket for both sending and receiving
        return True
    except:
        return False


def database_check(ip, port):
    """
    Test connection by opening a socket connection to url:port
    Retries int(i_r_c) times with int(i_w_d) secs wait between each retry
    """
    for i in range(1, int(i_r_c)):
        x = isOpen(ip, port)
        if x is False:
            isOpen(ip, port)
            print("fail", flush=True)
            # sys.stdout.flush()
        else:
            print("pass", x, flush=True)
            # sys.stdout.flush()
            break
    time.sleep(int(i_w_d))


def retryConnection(ip, port):
    """
    Retries int(l_r_c) times with int(l_t_s) secs wait between each retry
    """
    print("Retrying to establish connection", flush=True)
    for y in range(1, int(l_r_c)):
        z = isOpen(ip, port)
        if z is True:
            return z
        time.sleep(int(l_t_s))
    if z is False:
        return z


def liveness(ip, port):
    """
    Test connection by opening a socket connection to url:port
    Wait int(l_p_s) secs between each retry with max retry of int(l_r_c)
    """
    while True:
        res = isOpen(ip, port)
        if res is True:
            print("Liveness Running", flush=True)
            sys.stdout.flush()
        else:
            print("Liveness Failed", flush=True)
            sys.stdout.flush()
            z = retryConnection(ip, port)
            if z is False:
                print("Liveness finally failed:", flush=True)
                break
        time.sleep(int(l_p_s))


if __name__ == '__main__':
    url = sv+"."+ns+"."+"svc.cluster.local"
    ip = url
    port = int(port)
    database_check(ip, port)
    liveness(ip, port)
