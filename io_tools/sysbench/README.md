# Running the Sysbench benchmark on mysql/percona container/pod

-------------------------------------------------------------

## What is Sysbench benchmark

As an OLTP system benchmark, TPC-C (Transaction Performance Council - C) simulates a complete environment where a population of terminal operators executes transactions against a database. The benchmark is centered around the principal activities (transactions) of an order-entry environment. These transactions include entering and delivering orders, recording payments, checking the status of orders, and monitoring the level of stock at the warehouses.

## Steps to run Sysbench benchmark

- Obtain the IP of the DB container/pod. If it is a kubernetes pod, use the kubectl desribe command

  ```bash
  kubectl describe pod percona | grep IP
  ```

- Note the db_user name and password for the mysql to perform the remote login

- Pull the litmuschaos/sysbench docker image on the test host/kubernetes minion

  ```bash
  docker pull litmuschaos/sysbench
  ```

- Edit the sysbench.conf file to set the right values for the benchmark attributes

  ```json
  {
    "mysql-user": "root",
    "mysql-password": "password",
    "mysql-port": "3306",
    "db-driver": "mysql",
    "range_size": "100",
    "table_size": "10000",
    "tables": "2",
    "threads": "1",
    "events": "0",
    "time": "60",
    "rand-type": "uniform"
  }
  ```
  
- Run the sysbench-client container, bind-mounting the sysbench.conf file into the container as shown:
  ```bash
  docker run -it --link some-mysql:mysql -v $PWD:/tmp --rm sysbench mysql /tmp/sysbench.conf
  ```
