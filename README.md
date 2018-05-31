# Terraform vSphere OVF Provider
This is a *draft* [Terraform][terraform] provider that will allow users to upload 
OVF and OVA files to [vCenter][vmware-vcenter]. It is not ready for use yet!

[vmware-vcenter]: https://www.vmware.com/products/vcenter-server.html
[terraform]: https://github.com/hashicorp/terraform

## Current features
- upload an OVF template to vCenter
- mark the resulting VM as a template or virtual machine
- set network mappings on the template

## Upcoming features
- upload an OVA
- update the template
- export template UUID or name to be used by the [Terraform vSphere provider][provider]
- destroy the template
- otherwise customize the template
- documentation!

[provider]: https://github.com/terraform-providers/terraform-provider-vsphere