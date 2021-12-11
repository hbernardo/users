# users

## Requirements

- Golang 1.17.3
- Docker client 20.10.11
- Minikube 1.24.0
- kubernetes-cli (kubectl) client 1.23.0
- Helm 3.7.2

## Manual setup steps

### Building the docker image

It will test and build the application inside the Docker container.

```console
DOCKER_BUILDKIT=0 docker build \
    -t "hbernardo-users:$(git rev-parse --short HEAD)" \
    -t "hbernardo-users:latest" \
    -f go-src/Dockerfile .
```
