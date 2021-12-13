# Users API (HTTP Server)

This project implements a RESTful API in Golang that returns users information in JSON format.
The data is currently provided in a data file ["data/users.json"](data/users.json)).

The project also contains [Dockerfile](Dockerfile) and [Helm chart](helm-chart) for deploying it to Kubernetes cluster.

## Running the system locally

### Requirements (with tested versions)
- Docker client (20.10.11)
- Minikube (1.24.0)
- kubernetes-cli (kubectl) client (1.23.0)
- Helm (3.7.2)

### Usage

```console
./start.sh
```

The script ["start.sh"](start.sh) will setup the API locally using the port 8080 (http://localhost:8080).

It includes the following:
- Starting K8s cluster in Minikube
- Building the docker image (include running tests and build)
- Deploying the Helm Chart in the Minikube K8s cluster
- Accessing the service via K8s port-forwarding

Please check the [script content](start.sh) for manual execution of any step.

## Running only the application locally

### Requirements (with tested versions)
- Golang (1.17.3)

### Usage
```console
# Exporting required variables
export PORT=8080
export HEALTH_CHECK_PORT=8081
export LIVENESS_PROBE_PATH=/health/live
export READINESS_PROBE_PATH=/health/ready
export RATE_LIMIT_MAX_FREQUENCY=3
export RATE_LIMIT_BURST_SIZE=5
export RATE_LIMIT_MEMORY_DURATION=10m
export CORS_ALLOW_ORIGIN=http://localhost:8080
export CORS_ALLOW_METHODS=OPTIONS,GET,HEAD
export CORS_ALLOW_HEADERS=*
export LOG_LEVEL=debug

# Building the application
cd go-src
go build -o ../app ./cmd
cd -

# Running the HTTP server
./app http
```

## Exposed API routes

### GET users

Fetches multiple users based on pagination parameters ("limit" and "offset") got from the URL querystring.

#### Example:
[`http://localhost:8080/v1/users?limit=100&offset=100`](http://localhost:8080/v1/users?limit=100&offset=100)

### Path

`/users`

### Parameters and validations

- `limit` (querystring): required, positive integer less or equal to 1000
- `offset` (querystring): optional (default 0), positive integer

### Success response

  * **Code:** 200 <br/>
    **Content:** array of users data in JSON format

### Error response

  * **Code:** 500 (internal server error), 400 (bad request), 412 (precondition failed), 429 (too many requests), 304 (not modified) <br/>
    **Content:** `{"error": "{error information}"}`

### GET user by ID

Fetches single user by its ID (got from URL parameter).

#### Example:
[`http://localhost:8080/v1/users/f3f1612d-8239-4933-9891-71b5ee127844`](http://localhost:8080/v1/users/f3f1612d-8239-4933-9891-71b5ee127844)

### Path

`/users/{user_id}`

### Parameters

- `user_id` (url parameter): user ID (string)

### Success response

  * **Code:** 200 <br/>
    **Content:** user data in JSON format

### Error response

  * **Code:** 500 (internal server error), 404 (not found), 429 (too many requests), 304 (not modified) <br/>
    **Content:** `{"error": "{error information}"}`

## API structure design

### Command ["/cmd"](go-src/cmd) layer

Contains the application entrypoints and includes the ["main.go"](go-src/cmd/main.go) file.
I also instantiates/builds everything that is needed for the application execution.

### Server ["/srv"](go-src/srv) layer

Contains server implementation, including the HTTP server, handlers and middlewares.
This layer normally calls the "lib" layer implementations.

### Library ["/lib"](go-src/lib) layer

Contains business, domain and application logic implementation.
This layer normally calls the "infra" layer implementations.

### Infrastructure ["/infra"](go-src/infra) layer

Contains implementation related to data storing and processing or third-party integrations.
It's normally the lowest level part of the application.

## API middlewares

The API implements the following middlewares.

### ETag

Adds ETag header for proper client caching based on a version received as parameter (e.g. the users data version).

Returns 304 status code (not modified) if client requests the same version.

### Rate Limiter

RateLimiterMiddleware blocks the user from making a big amount of requests in a small amount of time.

It receives some configuration:
- Maximum allowed frequency (requests per second).
- Maximum bursts permitted.
- Duration of users rate limiter memory before it's cleaned.

### CORS

Sets proper [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS) headers to configure cross-origin access.

It also handles the preflight OPTIONS request.

### Panic Recovery

Treats any panic error that happens after this middleware and writes correct error to HTTP response and log.
