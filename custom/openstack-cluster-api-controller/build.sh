#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

if [ ! -f ${SCRIPT_DIR}/tmp/kubeadm ]; then
    mkdir -p ${SCRIPT_DIR}/tmp
    pushd ${SCRIPT_DIR}/tmp
    curl -L --remote-name-all https://storage.googleapis.com/kubernetes-release/release/v1.12.2/bin/linux/amd64/kubeadm
    popd
fi

docker build -t reg-dhc.app.corpintra.net/sbuerin/openstack-cluster-api-controller:latest .

docker --config ~/.docker-sbuerin push reg-dhc.app.corpintra.net/sbuerin/openstack-cluster-api-controller:latest
