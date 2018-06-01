---
name: Bug report
about: Create a report to help us improve

---

Hi there,

Thank you for opening an issue. Please note that we try to keep the Terraform
issue trackers reserved for bug reports and feature requests. For general usage
questions, please see: https://www.terraform.io/community.html or
https://communities.vmware.com/.

Also please note that this provider is a **draft** and is **unofficial**. Neither
Hashicorp nor VMWare currently provide any level of support for this Terraform 
provider. If you have issues using the provider that you think may be signs of
issues in Terraform or vSphere itself, please come to us first so we can make 
sure this is not our fault!

### Terraform Version

Run `terraform -v` to show the version. If you are not running the latest
version of Terraform, you can try upgrading, depending on if your problem is
with the provider or with Terraform core itself.

### Terraform Provider Versions

Run `terraform providers` to show the list of providers in use. If you are using the
official vSphere or NSX-T providers in your template, something like `provider.vsphere`
or `provider.nsxt` will show up in the output. If so, please report this version so that
we can investigate any compatibility issues between this provider and the vSphere
and NSX-T Terraform providers.

Since there is no official release of the vSphere OVF Terraform provider, you may need
to report the exact commit SHA that you built the OVF provider with.

### vSphere Version

There are some important differences between vSphere versions, especially in the way
that certain API calls involving VMs and VM templates are handled. Please provide the
version of vSphere that you are targeting with the provider.

### Affected Resource(s)

Please list the resources as a list, for example:
- `vsphere_virtual_machine`
- `vsphereovf_template`

If this issue appears to affect multiple resources, it may be an issue with
Terraform's core, so please mention this.

### Terraform Configuration Files

```hcl
# Copy-paste your Terraform configurations here - for large Terraform configs,
# please use a service like Dropbox and share a link to the ZIP file. For
# security, you can also encrypt the files using our GPG public key.
```

### OVA or OVF Templates

If you are using an OVA or OVF template that is publicly available, such as the
CoreOS OVA template, please provide a link to that template. If you are using
an OVA or OVF template that is not freely available to the general public but
is distributed as a supported product, such as the Pivotal Ops Manager OVA
template, please mention what product and version you are using. If you are
using a template that you built yourself, and are willing and able to provide
information on how to build this template, please do so, so that we may most
accurately recreate your problem.

### Debug Output

Please provide a link to a GitHub Gist containing the complete debug output:
https://www.terraform.io/docs/internals/debugging.html. Please do NOT paste the
debug output in the issue; just paste a link to the Gist.

### Panic Output

If Terraform produced a panic, please provide a link to a GitHub Gist containing
the output of the `crash.log`.

### Expected Behavior

What should have happened?

### Actual Behavior

What actually happened?

### Steps to Reproduce

Please list the steps required to reproduce the issue, for example:
1. `terraform apply`

### Important Factoids

Are there anything atypical about your infrastructure that we should know about
that could be causing an edge case or something not necessarily obvious? If so,
please state it here.

### References

Are there any other GitHub issues (open or closed) or Pull Requests that should
be linked here? For example:
- GH-1234
