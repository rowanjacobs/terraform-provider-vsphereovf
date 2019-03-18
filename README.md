# Terraform vSphere OVF Provider
This is a *draft* (not yet working!) [Terraform][terraform] provider that will
allow users to upload OVF and OVA files to [vCenter][vmware-vcenter]. It is not
ready for use yet!

**Note:**

So far, the provider can successfully upload and delete templates. However, it
cannot be used to feed the newly created template into the [Terraform vSphere provider][provider].

This will probably never actually work properly, barring some very unlikely
changes to the vSphere provider or the Terraform lifecycle. Here's why:

The `vsphere_virtual_machine` resource has a `CustomizeDiff` step that, when
cloning a template, [explicitly checks][vsphere-customize-diff] that template
exists as part of validation. Without having some way to short-circuit this
step, a user cannot upload a vSphere template and clone that template in the
same `terraform apply`.

Given that all the prospective users I have talked to have selected other means
of paving vSphere environments (either by abandoning Terraform entirely or
abandoning it in vSphere specifically), and given the (entirely reasonable)
resistance of the Terraform maintainers for supporting this use case, I do not
have very much external motivation to keep working on this project. Short of
discovering a way to short-circuit clone validation, this may be my last commit
to this repo.

[vmware-vcenter]: https://www.vmware.com/products/vcenter-server.html
[terraform]: https://github.com/hashicorp/terraform
[vsphere-customize-diff]: https://github.com/terraform-providers/terraform-provider-vsphere/blob/master/vsphere/resource_vsphere_virtual_machine.go#L726-L727
[provider]: https://github.com/terraform-providers/terraform-provider-vsphere

## What's done
- upload an OVF or OVA template to vCenter
- mark the resulting VM as a template or virtual machine
- set network mappings on the template
- export template UUID or name to be used by the vSphere provider
- destroy the template

## What's left to do
- figure out what's with the 4 minute gap between "creating template resource"
  and "importing template as vApp"
- find a way to bypass clone validation
- update the template
- otherwise customize the template
- documentation

