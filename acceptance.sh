#!/bin/bash
pushd vsphereovf > /dev/null
  TF_ACC=1 go test -v
popd
