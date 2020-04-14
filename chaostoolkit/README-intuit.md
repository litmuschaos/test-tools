# Chaos Toolkit for Litmus Chaos

# Local Development

* in this directory `cd test-tools/chaostoolkit`
* build python package `python setup.py develop`
* in this directory `cd test-tools`
* build pip module `pip install chaostoolkit/ `
* Initialize the submodule `git submodule init`
* Get the remote `git remote -v `
* Update the submodule `git submodule update`

# litmus
Litmus for chaos CRDS
* Get eiamCli Login `eiamCli login`
* export KUBECONFIG `export KUBECONFIG=/Users/snagal/.kube/admins@msaasfmea4-ppd-usw2.cluster.k8s.local`
* Apply operator - `kubectl apply -f litmus-operator-v1.0.0.yaml`
* Validate opertor - `kubectl get pods -n litmus`
* Verify CRDS - `kubectl get crds | grep chaos`
* Validate api-resource created - `kubectl api-resources | grep litmus`

# Base
* Apply experiments - `kubectl apply -f experiments.yaml`
* Validate - `kubectl get chaosexperiments`
* Annotate your app - `kubectl annotate deploy/fmea-test-4-appd-deployment litmuschaos.io/chaos="true"`
* Setup RBAC - please pick w2 for west - `kubectl apply -f rbac-pod-delete.yaml`
* Create Experiment - please pick w2 for west -`kubectl apply -f pod-delete.yaml`

# Intuit Changes
* Apply experiments for k8 - `kubectl apply -f experiments-k8.yaml`
* Validate the experiments for k8 - `kubectl get chaosexperiment`
* Apply experiments for aws - `kubectl apply -f experiments-aws.yaml`
* Validate the experiments for AWS - `kubectl get chaosexperiment`
* Setup RBAC - for pod delete RBAC - `kubectl apply -f rbac-pod-delete-app.yaml`
* Create pod Experiment - for count experiment -`kubectl create -f chaosengine-app-count.yaml`
* Create pod Experiment - for health experiment -`kubectl create -f chaosengine-app-health.yaml`
* Create pod Experiment - for health experiment -`kubectl create -f chaosengine-app-percentage.yaml`
* Validate the experiment -  `k describe chaosengine <k8-app-delete-health | k8-app-delete-count | k8-app-delete-percentage>`
* Execute the job - for debug pupose - `kubectl create -f k8-pod-delete-job.yaml`
 * Connect to pod and execute the `python3 k8_wrapper.py`

# Cleaning specific CR
* Delete the chaosengine ` kubectl delete chaosengine --all `
* Delete the jobs ` kubectl delete jobs --all `
* Delete the pod `kubectl delete pod <pod name>`

# Argo setup to see locally
* Download the artifact - `https://github.com/argoproj/argo/releases/tag/v2.5.0-rc11`
* Enable the argo Auth `~/argo-darwin-amd64 auth token`
* Bring the argo UI `argo-darwin-amd64 server -b -n argo`

# Litmus for chaos remove
* Delete Experiment - `kubectl delete -f pod-delete.yaml`
* Delete Pod - `kubectl delete pod <pod name>`
* Delete RBAC - `kubectl delete -f rbac-pod-delete.yaml`
* Delete experiments - `kubectl apply -f experiments.yaml`
* Validate api-resource created - `kubectl api-resources | grep litmus`
* Delete operator - `kubectl delete -f litmus-operator-v1.1.0.yaml`

# Litmus for intuit remove
* Delete Experiment - `kubectl delete -f chaosengine-app-health.yaml`
* Delete the chaosengine ` kubectl delete chaosengine --all `
* Delete the jobs ` kubectl delete jobs --all `
* Delete the pod `kubectl delete pod <pod name>`
* Delete RBAC - `kubectl delete -f rbac-pod-delete-app.yaml`
* Delete experiments - `kubectl delete -f experiments-k8.yaml`
* Validate api-resource created - `kubectl api-resources | grep litmus`
* Delete operator - `kubectl delete -f litmus-operator-v1.0.0.yaml`

# Litmus for local development
* Build the python project from  litmus/chaostoolkit-litmus -> `python setup.py develop`
* publish to the python library from litmus -> `pip install chaostoolkit-litmus/`
* Execute the python code for k8 `python k8_wrapper.py`
* Execute the python code for aws `python aws_wrapper.py`
