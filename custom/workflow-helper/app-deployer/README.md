## About App-Deployer
App-Deployer has been used for the installation of applications.

At first, the user is asked to give the namespace, typeName, operation, app, scope, and timeout.

<table>
  <tr>
    <th> Argument/Flag </th>
    <th> Description </th>
    <th> Example </th>
    <th> Domain </th>
  </tr>
  <tr>
    <td> -namespace </td>
    <td> Namespace of the application </td>
    <td> -namespace="sock-shop" </td>
    <td> sock-shop, loadtest, bank </td>
  </tr>
  <tr>
    <td> -typeName </td>
    <td> Type of application </br>
        weak scenario: It will create a single replica and Deployments for all</br> 
        resilient scenario: It will create two replicas of Statefulsets for databases </br> 
        and two replicas of Deployments for others.
    </td>
    <td> -typeName="resilient"</td>
    <td> resilient, weak</td>
  </tr>
  <tr>
    <td> -timeout </td>
    <td> Time after which it will fail the app-deployer if application is not ready </td>
    <td> -timeout="400" </td>
    <td> any integer </td>
  </tr>
  <tr>
    <td> -operation</td>
    <td> Operation will be Kubernetes CRUD operations on these resources. </br> 
         It supports create,delete, and apply operations.
    </td>
    <td> -operation="apply" </td>
    <td> create, apply, delete </td>
  </tr>
  </tr>
    <td> -app </td>
    <td> Name of the application </td>
    <td> -app="sock-shop" </td>
    <td> sock-shop, podtato-head, bank </td>
  </tr>
    <td> -scope </td>
    <td> Scope of application. It supports cluster and namespace modes </td>
    <td> -scope="cluster" </td>
    <td> cluster, namespace </td>
  </tr>
</table>
</br>

It creates a namespace and then installs the required application based on the given -namespace and -typeName.
If namespace already exists then it will skip the namespace creation and starts installing the application.
`[Status]: Namespace already exists!`

## Load-Test:
The load test packages a test script in a container for Locust that simulates user traffic to the application. Please run it against the front-end service. The address and port of the frontend will be different and depend on which platform you've deployed to. See the notes for each deployment.
It has been used parallelly with a chaos engine, which loads against the catalog front-end service.
In the manifest, it is written as:
```
- name: install-application
      container:
        image: litmuschaos/litmus-app-deployer:latest
        args: ["-namespace=loadtest"] 
```

## Now let’s see how the app-deployer works in the workflow:
At first, the installation of App-Deployer(application installation) is performed.
```
- name: install-application
      container:
        image: litmuschaos/litmus-app-deployer:latest
        args: ["-namespace=sock-shop","-typeName=resilient","-operation=apply","-timeout=400", "-app=sock-shop","-scope=cluster"] 
```
 ##### Note: 
  - By default, it will be resilient. For weak provide typeName flag as weak i.e,(-typeName=weak)
  - To delete application or loadtest corresponding operation need to pass `-operation=delete`

## Application details
As of now Sock-shop and Potato-Head applications have been added in pre-defined workflows for the weak and resilient cases.

### In a weak scenario, 
All the components of the application are deployments, in both sock-shop, podtato-head and bank-of-anthos. 
Single replica for every deployment.

#### For Databases:
MongoDB is run as a single replica deployment.
MySQL DB is hosted on ephemeral storage with a single replica.

Note: Only one replica of the deployment is present. After chaos injection, it will be down, and therefore service is unaccessible and eventually, it will fail due to resources unavailability.

### In a resilient scenario, 
Sock-shop: All the databases(MongoDB and MySQL) are statefulsets and the rest of components are deployments.
Podtato-head: All components are deployments.
bank-of-anthos: All the databases(MongoDB and MySQL) are statefulsets and the rest of components are deployments.
For each components(deployment/statefulset) two replicas for chaos pods are present. 

#### For Databases:
MongoDB multi-replica statefulset with persistent volumes. 
MySQL DB is Persistent, which can be dynamically extensible.
