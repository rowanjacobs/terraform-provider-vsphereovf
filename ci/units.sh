#!/usr/bin/env bash

pushd terraform-provider-vsphereovf/ > /dev/null
	ginkgo -r -p -v
popd > /dev/null
