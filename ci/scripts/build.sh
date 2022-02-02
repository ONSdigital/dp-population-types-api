#!/bin/bash -eux

pushd dp-population-types-api
  make build
  cp build/dp-population-types-api Dockerfile.concourse ../build
popd
