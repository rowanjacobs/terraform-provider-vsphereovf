#!/bin/bash
pushd vsphereovf > /dev/null
  TF_ACC="1" ginkgo -r -p -randomizeAllSpecs
popd
