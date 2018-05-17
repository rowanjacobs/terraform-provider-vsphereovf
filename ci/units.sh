#!/bin/bash -eux

mkdir -p "${GOPATH}/src/github.com/rowanjacobs"
ln -sf "${PWD}/terraform-provider-vsphereovf" "${GOPATH}/src/github.com/rowanjacobs"

pushd "${GOPATH}/src/github.com/rowanjacobs/terraform-provider-vsphereovf" > /dev/null
	ginkgo -r -p -v
popd > /dev/null
