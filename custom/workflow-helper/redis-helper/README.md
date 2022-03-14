## Load-Tested for Redis

Generating load against Redis using locust io. 

Locust helps in defining website user behavior with code and swarms your system with millions of simultaneous users.
Logs of the metrics contain requests per second, total requests, failed requests, average, min, max, and failed per second.

### To run load generator: 
 - Deploying load application and generator [Manifest](./manifest)
    - Redis application with single replica, `deployment.yaml`
    - Redis application with two replicas, `statefulset.yaml`
    - Redis load generate using locust, `redisLoadDeployment.yaml`
    - Redis load generate using script, `redisLoadLocustDeployment.yaml`
 - Load gen code is inside [Script](./load-gen)

### Get/Set requests are send to Redis Host

```
Name                                                                              # reqs      # fails  |     Avg     Min     Max  Median  |   req/s failures/s
----------------------------------------------------------------------------------------------------------------------------------------------------------------
 GET key0                                                                              12     0(0.00%)  |       0       0       1       0  |    2.00    0.00
 SET key0                                                                               4     0(0.00%)  |       1       0       7       0  |    0.50    0.00
 GET key1                                                                              12     0(0.00%)  |       0       0       1       0  |    2.00    0.00
 SET key1                                                                               4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key10                                                                             12     0(0.00%)  |       0       0       1       0  |    2.00    0.00
 SET key10                                                                              4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key100                                                                            12     0(0.00%)  |       0       0       0       0  |    2.00    0.00
 SET key100                                                                             4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key101                                                                            12     0(0.00%)  |       0       0       2       0  |    2.00    0.00
 SET key101                                                                             4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key102                                                                            12     0(0.00%)  |       0       0       0       0  |    2.00    0.00
 SET key102                                                                             4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key103                                                                            12     0(0.00%)  |       0       0       0       0  |    2.00    0.00
 SET key103                                                                             4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key104                                                                            12     0(0.00%)  |       0       0       0       0  |    2.00    0.00
 SET key104                                                                             4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key105                                                                            12     0(0.00%)  |       0       0       0       0  |    2.00    0.00
 SET key105                                                                             4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key106                                                                            12     0(0.00%)  |       0       0       0       0  |    2.00    0.00
 SET key106                                                                             4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key107                                                                            12     0(0.00%)  |       0       0       0       0  |    2.00    0.00
 SET key107                                                                             4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key108                                                                            12     0(0.00%)  |       0       0       0       0  |    2.00    0.00
 SET key108                                                                             4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key109                                                                            12     0(0.00%)  |       0       0       0       0  |    2.00    0.00
 SET key109                                                                             4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key11                                                                             12     0(0.00%)  |       0       0       0       0  |    2.00    0.00
 SET key11                                                                              4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key110                                                                            12     0(0.00%)  |       0       0       0       0  |    2.00    0.00
 SET key110                                                                             4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key111                                                                            12     0(0.00%)  |       0       0       0       0  |    1.50    0.00
 SET key111                                                                             4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key112                                                                            12     0(0.00%)  |       0       0       2       0  |    1.50    0.00
 SET key112                                                                             4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key113                                                                            12     0(0.00%)  |       0       0       1       0  |    1.50    0.00
 SET key113                                                                             4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key114                                                                            12     0(0.00%)  |       0       0       9       0  |    1.50    0.00
 SET key114                                                                             4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key115                                                                            12     0(0.00%)  |       0       0       0       0  |    1.50    0.00
 SET key115                                                                             4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
 GET key116                                                                            12     0(0.00%)  |       0       0       1       0  |    1.50    0.00
 SET key116                                                                             4     0(0.00%)  |       0       0       0       0  |    0.50    0.00
```
