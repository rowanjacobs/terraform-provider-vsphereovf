#!/bin/bash
pushd vsphereovf > /dev/null
  export TF_ACC="1"
  TF_ACC="1" go test -v -timeout 60m
popd
