# schema.py is used for getting secrets and helping read/write operations
import os
import base64
from kubernetes import client, config

if os.getenv('KUBERNETES_SERVICE_HOST'):
    configs = config.load_incluster_config()
else:
    configs = config.load_kube_config()

v1 = client.CoreV1Api()
namespace = os.getenv("NAMESPACE", "litmus")
secretName = os.getenv("SECRET_NAME", "postgres-application.credentials")

try:
    secret = v1.read_namespaced_secret(secretName, namespace)
except Exception as exp:
    print("[Error]: Unable to find secrets in litmus namespace err: ", exp)

# getSecret is returning details of postgres database
def getSecret():
    db_details = {}
    db_details["t_host"] = base64.b64decode(secret.data["host"]).decode('utf-8').replace('\n', '')
    db_details["t_port"] = base64.b64decode(secret.data["port"]).decode('utf-8').replace('\n', '')
    db_details["t_dbname"] = base64.b64decode(secret.data["dbname"]).decode('utf-8').replace('\n', '')
    db_details["t_name_user"] = base64.b64decode(secret.data["username"]).decode('utf-8').replace('\n', '')
    db_details["t_password"] =  base64.b64decode(secret.data["password"]).decode('utf-8').replace('\n', '')
    return db_details

# createSchema is creating a provided table in postgres database
def createSchema(t_name_tbl):
    s = "CREATE TABLE " + t_name_tbl + "( id serial NOT NULL"
    s += ", id_session int4 NULL DEFAULT 0, t_name_item varchar(64) NULL, t_contents text NULL, d_created date NULL DEFAULT now()"
    s += ", CONSTRAINT " + t_name_tbl + "_pkey PRIMARY KEY (id) ); "
    s += "CREATE UNIQUE INDEX " + t_name_tbl + "_id_idx ON public." + t_name_tbl + " USING btree (id);"
    return s

# dropSchema is dropping provided table
def dropSchema(t_name_tbl):
    return "DROP TABLE IF EXISTS " + t_name_tbl+";"
