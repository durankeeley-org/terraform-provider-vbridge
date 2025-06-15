---
page_title: "vbridge_virtual_machine Data Source - terraform-provider-vbridge"
subcategory: ""
description: |-
  This data source allows you to look up details about an existing virtual machine instance in Softsource vBridge.
---

# vbridge_virtual_machine (Data Source)

This data source allows you to retrieve metadata about an existing virtual machine instance in Softsource vBridge.

## Example Usage

```hcl
locals {
  subscription = "prod"
  christchurch_hosting = {
    christchurch_shortname       = "chch"
    christchurch_location        = "vcchcres"
    christchurch_locationname    = "Christchurch"
    christchurch_network         = "vcchcnet-prod"
  }
}

data "vbridge_virtual_machine" "example" {
  name      = "${local.subscription}-${local.christchurch_hosting.christchurch_shortname}-example"
  client_id = var.client_id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the virtual machine to look up.
* `client_id` - (Required) The ID of the vBridge client to which the VM belongs.

## Attributes Reference

The following attributes are exported:

* `client_id` - The ID of the client the VM belongs to.
* `name` - The name of the virtual machine.
* `guest_os_id` - The guest OS identifier.
* `cores` - Number of CPU cores allocated.
* `memory_size` - Amount of memory allocated in GB.
* `backup_type` - Backup configuration (`vBackupDisk`, `vBackupNone`).
* `hosting_location_id` - Hosting location ID.
* `operating_system_disk_guid` - GUID of the OS disk.
* `operating_system_disk_capacity` - Capacity of the OS disk in GB.
* `operating_system_disk_storage_profile` - Storage tier (e.g., `vStorageT1`).
* `vm_id` - Unique identifier of the virtual machine.
* `mo_ref` - The managed object reference (VMware).

## Notes

This data source is useful when you need to reference existing VMs in Terraform without managing their lifecycle. For example, if the VM is created externally or via another workspace, this allows you to reference its configuration and metadata.
