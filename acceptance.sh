#!/bin/bash
SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

COREOS_OVA_PATH="${SCRIPT_ROOT}/ignored/coreos_production_vmware_ova.ova"
COREOS_OVF_PATH="${SCRIPT_ROOT}/ignored/coreos_production_vmware_ova.ovf"

export VSPHERE_OVA_PATH=${VSPHERE_OVA_PATH:=${COREOS_OVA_PATH}}
export VSPHERE_OVF_PATH=${VSPHERE_OVF_PATH:=${COREOS_OVF_PATH}}
export VSPHERE_FOLDER=${VSPHERE_FOLDER:=khaleesi_templates}
export VSPHERE_DATACENTER=${VSPHERE_DATACENTER:=pizza-boxes-dc}
export VSPHERE_RESOURCE_POOL=${VSPHERE_RESOURCE_POOL:=khaleesi}
export VSPHERE_DATASTORE=${VSPHERE_DATSTORE:=vnx5600-pizza-2}
export VSPHERE_NETWORK=${VSPHERE_NETWORK:=khaleesi}

pushd vsphereovf > /dev/null
  TF_ACC="1" go test -v -timeout=12h
popd
