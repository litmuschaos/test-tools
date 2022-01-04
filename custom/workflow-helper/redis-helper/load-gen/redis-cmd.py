from pyredis import Client
import time, string, random, os

# DBDetails is performing io operations on redis db
class DBDetails(object):
    def __init__(self):
        self.client = Client(host=os.getenv("REDIS_HOST", "redis.redis.svc"))
        
    # createTable is inserting table in postgres database
    def createTable(self, client):

        randomstr = ''.join(random.choices(string.ascii_letters+string.digits,k=6))

        client.bulk_start()
        client.set(randomstr, randomstr, 10000)
        client.bulk_stop()

# main method  
def Main():
    try:
        db = DBDetails()
        db.createTable(db.client)
    except Exception as exp:
        print("[Error]: Unable to configure database err: {}".format(exp))
    print("0")

Main()