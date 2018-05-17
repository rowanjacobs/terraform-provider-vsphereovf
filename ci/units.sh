#!/usr/bin/env bash -eu

pushd terraform-provider-vsphereovf/ > /dev/null
	ginkgo -r -p -v
popd > /dev/null
