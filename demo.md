## Status Quo

We have 3 services in this repository that will need to be deployed to k8s in order for us to have the full app working.

- frontend

  This service contains the user facing interface loads all the messages from the board
- backend

  This service is responsible for receiving the requests from the users, process them and after that save it in the
  storage
- storage

  This service is responsible for storing the boards messages

### How to get these deployed:

- storage

  This service contains a chart that we will use to deploy the service.

  Steps
    ```bash
    # Build docker image
    docker build -t ghcr.io/joaopapereira/demo-intro-carvel/storage:1.16.0 .
    # Push image to the registry
    docker push ghcr.io/joaopapereira/demo-intro-carvel/storage:1.16.0
    # Install service on the cluster
    helm install storage chart --set image.repository=ghcr.io/joaopapereira/demo-intro-carvel/storage --create-namespace --namespace step1
    # Check that the pods are running
    kubectl get pods
    ```
- backend

  This service contains some ytt and kbld configuration.

  Steps
  ```bash
  # Single command that will generate the configuration, build image and deploy to the cluster
  ytt -f config | kbld -f- -f kbld.yaml | kapp deploy -a backend -f- -y
  ```

- frontend

  This service

  ```bash
  # Build docker image
  pack build ghcr.io/joaopapereira/demo-intro-carvel/frontend --path . --publish 
  # make sure the configuration is correct
  # apply to the cluster
  kubectl create -f config/deployment.yaml
  ```

### Trying to do some changes on the deployment without changing the code/base config

In a case where we want to standardize the installation with just a set of tools and want to make the life
of the people that are deploying the applications easier we can make carvel the default behavior

- storage

  This service contains a chart that we will use to deploy the service.

  Steps
    ```bash
    # Install service on the cluster
    helm template storage storage/chart --create-namespace | ytt -f- -f config-step2/storage -f config-step2/schema.yaml -f config-step2/values.yaml | kbld -f- | kapp deploy -a storage -f- -y
    # Check that the pods are running
    kubectl get pods
    ```
- backend

  This service contains some ytt and kbld configuration.

  Steps
  ```bash
  pushd backend
  # Single command that will generate the configuration, build image and deploy to the cluster
  ytt -f config -f ../config-step2/backend -f ../config-step2/schema.yaml -f ../config-step2/values.yaml | kbld -f- -f kbld.yaml | kapp deploy -a backend -f- -y
  
  popd
  ```

- frontend

  This service

  ```bash
  # Build docker image
  ytt -f go-frontend/config -f config-step2/frontend -f config-step2/schema.yaml -f config-step2/values.yaml | kbld -f- | kapp deploy -a frontend -f- -y
  ```

### What if the idea is to have this installed in a cluster and updated when something changes?

In this steps we can add kapp-controller to do this job for us

Installing kapp-controller

```bash
kapp deploy -a kc -f https://github.com/carvel-dev/kapp-controller/releases/download/v0.45.0/release.yml
```
