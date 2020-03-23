from pymongo import MongoClient
import time
import os
from pprint import pprint
import sys


# Assigning the environment variables
# Time period (in sec) b/w retries for DB init check
i_w_d = os.environ['INIT_WAIT_DELAY']
# No of retries for DB init check
i_r_c = os.environ['INIT_RETRY_COUNT']
# Time period (in sec) b/w liveness checks
l_p_s = os.environ['LIVENESS_PERIOD_SECONDS']
# Time period (in sec) b/w retries for db_connect failure
l_t_s = os.environ['LIVENESS_TIMEOUT_SECONDS']
# No of retries after a db_connect failure before declaring liveness fail
l_r_c = os.environ['LIVENESS_RETRY_COUNT']
ns = os.environ["NAMESPACE"]  # Namespace in which mongo is running
sv = os.environ["SERVICE_NAME"]  # Service name of mongodb


# function for db_init check
def db_init_check(db):
    try:
        connections_dict = db.command("serverStatus")["connections"]
        pprint(connections_dict)
        sys.stdout.flush()
        return 1
    except Exception as e:
        print("connection lost")
        sys.stdout.flush()
        return 0


# checking the database status
def database_check(db):
    pprint(db)
    for i in range(1, int(i_r_c)):
        x = db_init_check(db)
        if x == 0:
            db_init_check
            print("fail")
            sys.stdout.flush()
        else:
            print("pass", x)
            sys.stdout.flush()
            break
        time.sleep(int(i_w_d))


def retry_connection(db):
    for y in range(1, int(l_r_c)):
        z = db_init_check(db)
        if z == 1:
            return z
        time.sleep(int(l_t_s))
    if z == 0:
        return z


# liveness check
def liveness_check(db):
    while True:
        try:
            serverStatusResult = db.command("serverStatus")
            print("liveness Running")
            sys.stdout.flush()
        except Exception as e:
            print("liveness Failed", e)
            sys.stdout.flush()
            z = retry_connection(db)
            if z == 0:
                break   
        time.sleep(int(l_p_s))


if __name__ == "__main__":

    url = "mongodb://"+sv+"."+ns+"."+"svc.cluster.local/mydb"
    client = MongoClient(url)
    db = client.admin
    database_check(db)
    liveness_check(db)
