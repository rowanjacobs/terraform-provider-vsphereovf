# Terraform vSphere OVF Provider
This is a *draft* (not yet working!) [Terraform][terraform] provider that will allow users to upload 
OVF and OVA files to [vCenter][vmware-vcenter]. It is not ready for use yet!

[vmware-vcenter]: https://www.vmware.com/products/vcenter-server.html
[terraform]: https://github.com/hashicorp/terraform

## Current features
- upload an OVF or OVA template to vCenter
- mark the resulting VM as a template or virtual machine
- set network mappings on the template
- export template UUID or name to be used by the [Terraform vSphere provider][provider]

## Upcoming features
- read the template from other resources
- update the template
- destroy the template
- otherwise customize the template
- documentation

[provider]: https://github.com/terraform-providers/terraform-provider-vsphere
