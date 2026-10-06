#!/bin/bash -eux

pushd dp-population-types-api
  pip install poetry
  make -C sdk/python install-dev
  make audit
popd