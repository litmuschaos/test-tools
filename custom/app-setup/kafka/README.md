# Setup kafka cluster Using Ansible
## Follow the steps to setup kafka cluster
### Step:1
- Create  k8s ServiceAccount\
 Use this sample RBAC manifest to create serviceAccount.\
 **Note:** This example have all the role permissions.You can change it to minimum  necessary role permission as per your requirement.
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
    name: litmus-kafka-sa
    namespace: default
    labels:
      name: litmus-kafka-sa
---

apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: litmus-kafka-sa
  labels:
    name: litmus-kafka-sa
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["*"]
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: litmus-kafka-sa
  labels:
    name: litmus-kafka-sa
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: litmus-kafka-sa
subjects:
- kind: ServiceAccount
  name: litmus-kafka-sa
  namespace: default
```
 ### Step 2: 
- Create a k8s-secret and provide your aws credentials If your k8s-cluster is on aws-eks
```yaml
apiVersion: v1
kind: Secret
metadata:
        name: aws-secret
data:
        AWS_ACCESS_KEY_ID: "your base64-encoded access key"   
        AWS_SECRET_ACCESS_KEY: "your base64-encoded secret key"
        AWS_DEFAULT_REGION: "your base64-encoded region"
        EKS_CLUSTER_NAME:  "your base64-encoded cluster name"
        
```
### Step 3:
- Create a litmus-kafka-deployer pod  that will setup your kafka cluster
 #### Supported ENV variables
 <table>
    <tr>
      <th> Variables </th>
      <th> Description </th>
      <th> Specify In Pod </th>
      <th> Notes </th>
  </tr>
  <tr>
    <td> MODE </td>
    <td> Provide <code>setup</code> to setup kafka cluster and <code>cleanup</code> for uninstallation </td>
    <td> Required </td>
    <td> It supports two value <br/>
         MODE: setup <br/>
         MODE: cleanup </td>
  </tr>
  <tr>
    <td> PLATFORM </td>
    <td>  Provide the target platform  name. Like <code>eks</code> </td>
    <td> Optional </td>
    <td> Currently it supports only eks cluster <br/>
      PLATFORM: eks</td>
  </tr>
  <tr>
    <td> KUDO_VERSION </td>
    <td> Provide the Kudo version for kafka Installation  </td>
    <td> Optional </td>
   <td> If KUDO_VERSION is not provided ,By-default It will Install the 0.12.0 version of  KUDO </td>
  </tr>
  <tr>
   <td> KAFKA_NAMESPACE </td>
   <td> Provide the namespace for kafka </td>
   <td> Required</td>
   <td> Example: <br/>
        KAFKA_NAMESPACE: kafka </td>
 </tr>
  <tr>
  <td> MONITORING_NAMESPACE </td>
  <td> Provide the namespace for prometheus and grafana </td>
  <td> Required </td>
  <td> Example: <br/>
   MONITORING_NAMESPACE: monitoring </td>
 </tr>
  <tr>
  <td> LITMUS_GIT_TAG </td>
  <td> Provide the tag of Litmus github Repository </td>
  <td> Required </td>
  <td> Example: <br/>
   LITMUS_GIT_TAG: 2.0.0-Beta7</td>
 </tr>
 </table>
Use this Example to create litmus-kafka-deployer-pod.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: litmus-kafka-deployer
  labels:
    app: litmus-kafka-deployer
spec:
  serviceAccountName: litmus-kafka-sa
  containers:
  - name: litmus-kafka-deployer-container
    image: litmuschaos/kafka-deployer:latest
    imagePullPolicy: Always
    envFrom:
        - secretRef:
              name: aws-secret
    env:
        ##  It defines the mode of the experiment
        ##Supported values: setup, cleanup
        - name: MODE
          value: "setup"
          
         ## It defines the platform of the k8s cluster
         ## Supported value: eks
        - name: PLATFORM
          value: "eks"
        
        ## It defines the kubectl-Kudo version
        - name: KUDO_VERSION
          value: 0.12.0
          
        ## Name of namespace where kafka will be installed
        - name: KAFKA_NAMESPACE
          value: "kafka"

        ## Name of namespace where prometheus and grafana will be installed
        - name: MONITORING_NAMESPACE
          value: "monitoring"

        ## Tag of LITMUS Repository that will be clone
        - name: LITMUS_GIT_TAG
          value: 2.0.0-Beta7

 
```

It takes few minutes to setup the kafka cluster

