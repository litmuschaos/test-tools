from pyredis import Client
import time, string, random, os

host = os.getenv("REDIS_HOST", "redis.redis.svc")

# DBDetails is performing io operations on redis db
class DBDetails(object):
    def __init__(self):
        self.client = Client(host=host)
        
    # createTable is inserting table in postgres database
    def createTable(self, client):

        randomstr = ''.join(random.choices(string.ascii_letters+string.digits,k=6))
        
        client.bulk_start()
        client.set(randomstr, randomstr)
        
        client.bulk_stop()
        print("Added ", randomstr)

# main method
def Main():
    print("[Info]: Redis Load-generation has been started!!")
    print("[Info]: Redis dbname: {}, host: {}, port: {}".format(0, host, 6379))
    while True:
        try:
            db = DBDetails()
            db.createTable(db.client)
        except Exception as exp:
            print("[Error]: Unable to configure database err: {}".format(exp))
            time.sleep(1)
            continue
        time.sleep(0.5) 
Main()