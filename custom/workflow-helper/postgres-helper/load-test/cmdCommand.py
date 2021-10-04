import psycopg2 
import time, random
import string
from schema import getSecret, createSchema, dropSchema

# getSecret returns the db details
db_details = getSecret()

# DBDetails is connecting to a specific database and performing read and write operations
class DBDetails(object):
    def __init__(self, db_details):
        self.db_conn      = psycopg2.connect(host=db_details["t_host"], port=db_details["t_port"], dbname=db_details["t_dbname"], user=db_details["t_name_user"], password=db_details["t_password"])
        self.db_cursor    = self.db_conn.cursor()

    # createTable is inserting table in postgres database
    def createTable(self, db):

        db.db_cursor = db.db_conn.cursor()
        letters = string.ascii_lowercase
        t_name_tbl = ''.join(random.choice(letters) for i in range(5)) 
        s = createSchema(t_name_tbl)
        db.db_cursor.execute(s)
        db.db_conn.commit()

        return t_name_tbl

    # dropTable is deleting table from postgres database   
    def dropTable(self, db, t_name_tbl):

        s = dropSchema(t_name_tbl)
        db.db_cursor.execute(s)
        db.db_conn.commit()

# Main method of cmdCommand.py      
def Main():
    db = DBDetails(db_details)
    table = db.createTable(db)
    db.dropTable(db, table)
    print(0)

Main()
