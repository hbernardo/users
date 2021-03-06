# Test and build
FROM golang:1.17.3 as builder

WORKDIR /build

COPY ./go-src /build
COPY ./data /build/data

# running the tests
RUN go test -p 1 -v -race ./...
# building the application
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -a -v -o app ./cmd


# Run server
FROM busybox:latest
COPY --from=builder /build/app /users
COPY --from=builder /build/data /data

# default command
CMD ["/users", "http"]
