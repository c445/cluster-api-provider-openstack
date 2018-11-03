
docker pull gcr.io/k8s-cluster-api/cluster-api-controller:latest

docker tag gcr.io/k8s-cluster-api/cluster-api-controller:latest reg-dhc.app.corpintra.net/sbuerin/cluster-api-controller:latest

docker --config ~/.docker-sbuerin push reg-dhc.app.corpintra.net/sbuerin/cluster-api-controller:latest
