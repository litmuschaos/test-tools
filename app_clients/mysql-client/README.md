## litmuschaos/mysql-client

This image provides a couple of simple bash scripts that can be used alongside percona/mysql deployments to test their health & functionality

**Loadgen (MySQLLoadGenerate.sh)**: This script creates sample database with some initial data and inserts tables in an exponential manner 
(`INSERT INTO <table> select * FROM <table>` ).It takes the IP address of the mysql service as the argument & can be used in sample loadgen jobs.

**Liveness (mysql-liveness-check.sh)**: This script checks the status of the mysql service periodically with pre-defined interval & timeouts.
It takes a JSON file with database credentials as argument & can be used in external liveness jobs such as the one defined [here](mysql-liveness-check.yaml), with the credentials passed as a configmap, as shown [here](db-cred.cnf)

*TODO:Reuse the db credentials json in the loadgen too*







