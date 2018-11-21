# openshift-applier-operator

Applying OpenShift configurations in a declarative and automated fashion using the [openshift-applier]

## Custom Resources

A _CustomResourceDefinition_ is used to define a new object called `Applier` which specifies which git repository contains openshift-applier inventories along with how they can be triggered. 

The following sections describe the primary components of the _Applier_ resource

### Retrieving Source Code

Inventories leveraged by the OpenShift applier can be retrieved from Git based repositories which are defined within the 'Git" section underneath `source` within the `Applier` custom resource as shown below:

```
apiVersion: cop.redhat.com/v1alpha1
kind: Applier
metadata:
  name: applier-example
spec:
  source:
    git:
      uri: https://github.com/redhat-cop/openshift-applier
      inventoryDir: tests/inventories/params-from-file
```

The example below demonstrates how to specify the location of the repository and the relative path of where the inventory resources are located

#### Private Repositories

Private repositories leveraging the SSH protocol can be utilized by specifying the SSH based address of the repository and creating a secret containing the ssh key.

Create a secret containing the SSH key using the following command:

```
oc create secret generic applier-ssh --from-file=id_rsa=<LOCATION_OF_PRIVATE_KEY>
```

Add the `secretName` property underneath the `git` section as shown below:

```
apiVersion: cop.redhat.com/v1alpha1
kind: Applier
metadata:
  name: applier-example
spec:
  source:
    git:
      uri: git@github.com:redhat-cop/openshift-applier.git
      inventoryDir: tests/inventories/params-from-file
      secretName: applier-ssh
```

### Triggering Execution

The following methods are available for triggering the execution:

#### Webhook

A webhook (or HTTP based invocation) is available which uses a token to protect the endpoint. The `token` field of the `Applier` custom resource specifies the value that can be used at invocation.

The following is an example of how the token can be defined on the custom resource:

```
apiVersion: cop.redhat.com/v1alpha1
kind: Applier
metadata:
  name: applier-example
  namespace: applier-project
spec:
  source:
    git:
      ...
  webhook:
    token: secrettoken
```

The webhook can be invoked by submitting a post request to the following location:

```
<LOCATION_OF APPLICATION>/webhook/<NAMESPACE>/<TOKEN>
```

A 201 HTTP response will 

### Service Account

The execution of the applier operates using a service account. By default, service accounts do not have rights to query the majority of the API. Permissions must be given to this service account in the namespace containing the `Applier` resources. Alternatively, the `serviceAccount` can be used to specify another service account to utilize instead  


## Quickstart

Use the following steps to demonstrate the functionality of the project (Note: These steps should be executed as a user with `cluster-admin` privileges):

* Clone Repository

```
mkdir -p $GOPATH/src/github.com/redhat-cop
cd $GOPATH/src/github.com/redhat-cop
git clone <GIT_URL>
cd openshift-applier-operator
```

* Deploy the custom resource

```
oc apply -f deploy/crds/crd.yml
```

* Create a new Project

```
oc new-project applier-operator
```

* Create a new service account to run the applier job pod

```
oc create serviceaccount applier
```

* Grant the service account some permissions

```
oc adm policy add-cluster-role-to-user self-provisioner -z applier
```

* Create the example `Applier` resource

```
oc apply -f examples/applier-example.yml
```

* Start up the project using the operator-sdk (running locally)

```
operator-sdk up local
```

* Send a webhook post request to the operator

```
curl -X POST http://localhost:8080/webhook/applier-operator/securetoken
```

The pods executing the openshift-applier are then launched provisioning resources within the environment
