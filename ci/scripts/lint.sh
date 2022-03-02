#!/bin/bash -eux

pushd dp-population-types-api
  go install github.com/golangci/golangci-lint/cmd/golangci-lint
  go build
  ls -la
  make lint
popd