# **Git-App-Deployer**

Git-App-Deployer has been used for installation of sock-shop applications.
At first the user is asked to give the namespace, filePath and timeout.

## For -namespace :

Namespace provides an additional qualification to a unique resource name. This is helpful when multiple teams are using the same cluster and there is a potential of name collision. It can be a virtual wall between multiple clusters.

For "sock-shop" user has to pass
- -namespace=sock-shop

For Load-Test

- -namespace=loadtest

## For -filePath

For Weak Sock-Shop-Resilient check just pass:

- -namepsace=weak

In a weak scenario it will create single replica and Deployments for all.

For Resilient Sock-Shop-Resilient check just pass:

- -namepsace=resilient

In a resilient scenario it will create 2 replicas of pods with Statefulset for databases and Deployments for others.

## For -timeout

Timeout is used for termination of application. The exceeding time by default is 300s.
You may change the default time value e.g

- -timeout=400

A kubeconfig file is a file used to configure access to Kubernetes when used in conjunction with the kubectl command line tool (or other clients).
It creates a namespace and then installs the required application on the basis of the given -namespace and -filepath.
If namespace is already in exist then it shows log and start installing sock-shop.

 - [Status]: Namespace already exist!

Sock-shop installation will basically deploy all 14 manifest of sock-shop microservices.

At first the installation of Git-App-Deployer(Application installation) is done .
```
- name: install-application
      container:
        image: litmuschaos/litmus-app-deployer:latest
        args: ["-namespace=sock-shop","-typeName=weak", "-timeout=400"] 
```

Note :
- for resilient provide type flagName as resilient(-typeName=resilient)

In weak scenario single replica will be there. You may check using this command:
- kubectl get po -n sock-shop


In the terminal the output will be weak

```
NAME                            READY   STATUS    RESTARTS   AGE
carts-754c96bf74-hbwlv          1/1     Running   0          4m8s
carts-db-69f74c65bd-xkm8s       1/1     Running   0          4m9s
catalogue-67f86d6587-8gghr      1/1     Running   0          4m8s
catalogue-db-66bc8df878-km7m9   1/1     Running   0          4m8s
front-end-5d96ff485b-wkn75      1/1     Running   0          4m8s
orders-67dbccdfdf-bz4sp         1/1     Running   0          4m8s
orders-db-79c689d8c7-r52qc      1/1     Running   0          4m8s
payment-6787495757-p7l6l        1/1     Running   0          4m8s
queue-master-67478795b7-z8wxz   1/1     Running   0          4m8s
rabbitmq-c8fcd79c9-gzdsc        1/1     Running   0          4m7s
shipping-d49997689-zm72c        1/1     Running   0          4m7s
user-7498444df6-l8wqm           1/1     Running   0          4m6s
user-db-64b9d4f4d-2z6zj         1/1     Running   0          4m7s
user-load-dc4586796-5vwvz       1/1     Running   0          4m7s
```


In the terminal 2 replicas will be shown for resilient.

```
NAME                            READY   STATUS    RESTARTS   AGE
carts-754c96bf74-hbwlv          1/1     Running   0          19m
carts-754c96bf74-w7bzt          1/1     Running   0          3m48s
carts-db-0                      1/1     Running   0          3m48s
carts-db-1                      1/1     Running   0          3m44s
carts-db-69f74c65bd-xkm8s       1/1     Running   0          19m
catalogue-67f86d6587-k4ncs      1/1     Running   0          3m48s
catalogue-67f86d6587-r4mpl      1/1     Running   0          12m
catalogue-db-0                  1/1     Running   0          3m48s
catalogue-db-1                  1/1     Running   0          3m29s
catalogue-db-66bc8df878-km7m9   1/1     Running   0          19m
front-end-5d96ff485b-lw7vp      1/1     Running   0          3m47s
front-end-5d96ff485b-wkn75      1/1     Running   1          19m
orders-67dbccdfdf-bz4sp         1/1     Running   0          19m
orders-67dbccdfdf-rxll5         0/1     Running   0          3m47s
orders-db-0                     1/1     Running   0          3m47s
orders-db-1                     1/1     Running   0          3m33s
orders-db-79c689d8c7-r52qc      1/1     Running   0          19m
payment-6787495757-7njtm        1/1     Running   0          3m47s
payment-6787495757-p7l6l        1/1     Running   0          19m
queue-master-67478795b7-76n2n   1/1     Running   0          3m47s
queue-master-67478795b7-z8wxz   1/1     Running   0          19m
rabbitmq-c8fcd79c9-2sjxr        1/1     Running   0          3m47s
rabbitmq-c8fcd79c9-gzdsc        1/1     Running   0          19m
shipping-d49997689-l9gls        1/1     Running   0          3m47s
shipping-d49997689-zm72c        1/1     Running   0          19m
user-7498444df6-6t8cq           1/1     Running   0          3m47s
user-7498444df6-l8wqm           1/1     Running   0          19m
user-db-0                       1/1     Running   0          3m47s
user-db-1                       1/1     Running   0          3m1s
user-db-64b9d4f4d-2z6zj         1/1     Running   0          19m
user-load-dc4586796-5vwvz       1/1     Running   0          19m
user-load-dc4586796-wsh59       1/1     Running   0          3m47s
```

# **Load-Test**:
The load test packages a test script in a container for Locust that simulates user traffic to Sock Shop, please run it against the front-end service. The address and port of the frontend will be different and depend on which platform you've deployed to. See the notes for each deployment.
It has been used parallely with a chaos engine which loads against the catalogue front-end service.

In manifest it is written as :
```
- name: install-application
      container:
        image: litmuschaos/litmus-app-deployer:latest
        args: ["-namespace=loadtest"] 
```

- Load-test have 2 replicas as shown below :
```
oumkale@mayadata:~$ kubectl get po -n loadtest
NAME                         READY   STATUS    RESTARTS   AGE
load-test-5d489d8c9d-mxc5g   1/1     Running   0          88s
load-test-5d489d8c9d-qnrbs   1/1     Running   0          88s
```