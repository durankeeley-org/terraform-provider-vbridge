---
page_title: "vbridge_virtual_machine_additionaldisk Resource - terraform-provider-vbridge"
subcategory: ""
description: |-
  This resource allows you to add virtual machine additional disks in vBridge.
---

# vbridge_virtual_machine_additionaldisk (Resource)

This resource allows you to add additional disks to virtual machine instances in vBridge.

## Example Usage

```terraform
# Performance Disk
resource "vbridge_virtual_machine_additionaldisk" "disk2" {
  vm_id = resource.vbridge_virtual_machine.example.vm_id
  storage_profile = "vStorageT1"
  capacity = 47
}
```

## Argument Reference

The following arguments are supported:

* `vm_id` - (Required) The ID of the virtual machine to attach the disk to.

* `storage_profile` - (Required) The storage profile to use for the disk.

* `capacity` - (Required) The capacity of the disk in GB.
