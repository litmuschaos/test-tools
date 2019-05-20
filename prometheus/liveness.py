"""This script checks liveness of prometheus service by trying to establish socket connection using ip:port pair"""
import os
import socket
import sys
import time
import subprocess

''' Fetching below environment variables from litmus job  '''

# Number of retries to check liveness
LIVENESS_RETRY = os.environ['LIVENESS_RETRY_COUNT']

# Time period (in sec) between retries for liveness check
LIVENESS_TIMEOUT = os.environ['LIVENESS_TIMEOUT_SECONDS']

# Time period (in sec) between liveness checks
LIVENESS_PERIOD = os.environ['LIVENESS_PERIOD_SECONDS']

# Namespace where app is running
NAMESPACE = os.environ['APP_NAMESPACE']

# Service endpoint 
SERVICE_ENDPOINT = os.environ['SERVICE_ENDPOINT']

# Service port 
PORT = os.environ['PORT']

def is_open(ip, port):
    """Socket connection to url:port"""
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        s.connect((ip, int(port)))
        s.shutdown(2)  # Disables the socket for both sending and receiving
        return True
    except:
        return False

def retry_connection(ip, port):
    """Retries int(LIVENESS_RETRY) times with int(LIVENESS_TIMEOUT) secs wait between each retry"""
    print("Retrying to establish connection", flush=True)
    for itr in range(1, int(LIVENESS_RETRY)):
        result = is_open(ip, port)
        if result is True:
            return result
        time.sleep(int(LIVENESS_TIMEOUT))
    if result is False:
        return result

def liveness(ip, port):
    """Test connection by opening a socket connection to url:port 
    Wait int(LIVENESS_PERIOD) secs between each retry with max retry of int(LIVENESS_RETRY)"""
    while True:
        res = is_open(ip, port)
        if res is True:
            print("Liveness Running", flush=True)
            sys.stdout.flush()
        else:
            print("Liveness Failed", flush=True)
            sys.stdout.flush()
            result = retry_connection(ip, port)
            if result is False:
                print("Liveness finally failed", flush=True)
                break
        time.sleep(int(LIVENESS_PERIOD))

def main():
    ''' Checking app liveness by trying to establish socket connection to service endpoint '''
    print("Service Endpoint: {}".format(SERVICE_ENDPOINT))
    print("Port: {}".format(PORT))
    liveness(SERVICE_ENDPOINT, PORT)

main()
