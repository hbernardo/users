#!/usr/bin/env bash
set -e

COMMIT_HASH=$(git rev-parse --short HEAD)

DOCKER_BUILDKIT=0 docker build \
    -t "hbernardo-users:$(git rev-parse --short HEAD)" \
    -t "hbernardo-users:latest" \
    -f go-src/Dockerfile .
