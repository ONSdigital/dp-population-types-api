#!/bin/bash -eux

pushd dp-population-types-api
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2
  go build
  ls -la
  make lint
popd
