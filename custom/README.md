
# Installation @ CaaS

## Prerequisites

* create c01p006 tenant 
* create c01p008 tenant and delete all vms
````
dhc_openstack c01p008
openstack server list -f value -c ID | xargs openstack server delete
````
* get Secrets
````
custom/copy-secrets.sh
````

## Create cluster c01p008 with bootstrap cluster c01p006

### Update bootstrap cluster kubeconfig
````
dhc_kubeconfig c01p006 $GOPATH/src/sigs.k8s.io/cluster-api-provider-openstack/custom/secrets/c01p008/bootstrap.kubeconfig
kcfg $GOPATH/src/sigs.k8s.io/cluster-api-provider-openstack/custom/secrets/c01p008/bootstrap.kubeconfig
kctx c01p006-cluster-admin
````

### Run IDE config for cluster-api-provider-openstack c01p008, cmd/manager/main.go with the following env vars
````
KUBECONFIG=custom/secrets/c01p008/bootstrap.kubeconfig
OS_CLOUD=openstack
OS_CLIENT_CONFIG_FILE=custom/secrets/c01p008/clouds.yaml
SSH_KEY_USER_PATH=custom/secrets/user
SSH_PUBLIC_KEY_PATH=custom/secrets/public
SSH_PRIVATE_KEY_PATH=custom/secrets/private
MACHINE_SETUP_CONFIG_PATH=custom/config/machine_setup_configs.yaml
````

### Run IDE clusterctl config c01p008 or execute the following
````
PROJECT_DIR=$GOPATH/src/sigs.k8s.io/cluster-api-provider-openstack
cd $PROJECT_DIR/cmd/clusterctl
./clusterctl create cluster --existing-bootstrap-cluster-kubeconfig $PROJECT_DIR/custom/secrets/c01p008/bootstrap.kubeconfig --provider openstack -c $PROJECT_DIR/custom/config/c01p008/cluster.yaml -m $PROJECT_DIR/custom/config/c01p008/machines.yaml -p $PROJECT_DIR/custom/secrets/c01p008/provider-components.yaml --cleanup-bootstrap-cluster=true
k get cluster,machine
````

## Create cluster c01p006 with bootstrap cluster c01p008

### Update bootstrap cluster kubeconfig
````
dhc_kubeconfig c01p008 $GOPATH/src/sigs.k8s.io/cluster-api-provider-openstack/custom/secrets/c01p006/bootstrap.kubeconfig
kcfg $GOPATH/src/sigs.k8s.io/cluster-api-provider-openstack/custom/secrets/c01p006/bootstrap.kubeconfig
kctx c01p008-cluster-admin
````

### Run IDE config for cluster-api-provider-openstack c01p006, cmd/manager/main.go with the following env vars
````
KUBECONFIG=custom/secrets/c01p006/bootstrap.kubeconfig
OS_CLOUD=openstack
OS_CLIENT_CONFIG_FILE=custom/secrets/c01p008/clouds.yaml
SSH_KEY_USER_PATH=custom/secrets/user
SSH_PUBLIC_KEY_PATH=custom/secrets/public
SSH_PRIVATE_KEY_PATH=custom/secrets/private
MACHINE_SETUP_CONFIG_PATH=custom/config/machine_setup_configs.yaml
````

### Run IDE clusterctl config c01p006 or execute the following
````
PROJECT_DIR=$GOPATH/src/sigs.k8s.io/cluster-api-provider-openstack
cd $PROJECT_DIR/cmd/clusterctl
./clusterctl create cluster --existing-bootstrap-cluster-kubeconfig $PROJECT_DIR/custom/secrets/c01p006/bootstrap.kubeconfig --provider openstack -c $PROJECT_DIR/custom/config/c01p006/cluster.yaml -m $PROJECT_DIR/custom/config/c01p006/machines.yaml -p $PROJECT_DIR/custom/secrets/c01p006/provider-components.yaml --cleanup-bootstrap-cluster=true
k get cluster,machine
````

# Open Issues

* floating ip per machine
* network id
* security group?
* controller needs hostpath mounts
  * kubeadm
  * certs
  * config
  
  
# Additional Info

## Build clusterctl
````
git clone https://github.com/kubernetes-sigs/cluster-api-provider-openstack $GOPATH/src/sigs.k8s.io/cluster-api-provider-openstack
cd $GOPATH/src/sigs.k8s.io/cluster-api-provider-openstack/cmd/clusterctl
go build
````

## Generate YAMLs

````
cd examples/openstack
./generate-yaml.sh
# insert config of an openstack tenant (e.g. c01p006)
cd ../..
````
  