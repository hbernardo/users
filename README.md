# Users API (HTTP Server)

## Running on local machine

### Requirements (with tested versions)
- Bash (for executing the shell script)
- Docker client (20.10.11)
- Git (2.25.0)
- Golang (1.17.3)
- Minikube (1.24.0)
- kubernetes-cli (kubectl) client (1.23.0)
- Helm (3.7.2)

### Single command

```console
./start.sh
```

### Manual setup steps (executed by the script "start.sh")

#### Starting K8s cluster in Minikube

```console
minikube start
```

#### Building the docker image

It will test and build the application inside the Docker container.

```console
eval $(minikube docker-env) # setting minikube docker env

DOCKER_BUILDKIT=0 docker build \
    -t "hbernardo-users:$(git rev-parse --short HEAD)" \
    -t "hbernardo-users:latest" \
    -f go-src/Dockerfile .
```

#### Deploying the Helm Chart

```console
helm upgrade --install users-api-http-server -f helm-chart/values.yaml helm-chart
kubectl rollout restart deployment users-api-http-server
kubectl rollout status deployment users-api-http-server
```

#### Accessing the service via K8s port-forwarding

```console
export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=users,app.kubernetes.io/instance=users-api-http-server" -o jsonpath="{.items[0].metadata.name}")
export CONTAINER_PORT=$(kubectl get pod --namespace default $POD_NAME -o jsonpath="{.spec.containers[0].ports[0].containerPort}")
echo "Visit http://127.0.0.1:8080 to use your application\n"
kubectl logs $POD_NAME -f & # follow the logs
kubectl --namespace default port-forward $POD_NAME 8080:$CONTAINER_PORT
```
