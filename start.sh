#!/usr/bin/env bash
set -e

printf "# Starting K8s cluster in Minikube\n"

minikube start


printf "\n# Building the docker image (test and build the application)\n"

eval $(minikube docker-env) # setting minikube docker env
DOCKER_BUILDKIT=0 docker build \
    -t "hbernardo-users:$(git rev-parse --short HEAD)" \
    -t "hbernardo-users:latest" \
    -f go-src/Dockerfile .


printf "\n# Deploying the Helm Chart\n"

helm upgrade --install users-api-http-server -f helm-chart/values.yaml helm-chart
kubectl rollout restart deployment users-api-http-server
kubectl rollout status deployment users-api-http-server


printf "\n# Accessing the service via K8s port-forwarding\n"

sleep 10 # making sure everything is updated inside the cluster
export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=users,app.kubernetes.io/instance=users-api-http-server" -o jsonpath="{.items[0].metadata.name}")
export CONTAINER_PORT=$(kubectl get pod --namespace default $POD_NAME -o jsonpath="{.spec.containers[0].ports[0].containerPort}")
printf "Visit http://127.0.0.1:8080 to use your application"
kubectl --namespace default port-forward $POD_NAME 8080:$CONTAINER_PORT
