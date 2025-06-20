---
page_title: "vbridge_virtual_machine Resource - terraform-provider-vbridge"
subcategory: ""
description: |-
  Provides a vBridge virtual machine resource for managing virtual machine instances on the Softsource vBridge platform.
---

# vbridge_virtual_machine (Resource)

The `vbridge_virtual_machine` resource allows you to create and manage virtual machines in [Softsource vBridge](https://www.svbgroup.co.nz/), a New Zealand-based public cloud provider headquartered in Christchurch.

## Example Usage

```hcl
resource "vbridge_virtual_machine" "example" {
  client_id   = var.client_id
  name        = "prod-chch-web01"
  template    = "Windows2022_Standard_30GB"
  guest_os_id = "windows2019srv_64Guest"
  cores       = 2
  memory_size = 6

  operating_system_disk_capacity        = 35
  operating_system_disk_storage_profile = "vStorageT1"

  hosting_location_id            = "vcchcres"
  hosting_location_name          = "Christchurch"
  hosting_location_default_network = "vcchcnet-prod"

  backup_type = "vBackupDisk"

  tags        = { environment = "prod" }
  description = "Web frontend"
  notes       = "Created via Terraform"
}
````

## Create Logic Explained

The creation of a `vbridge_virtual_machine` resource follows these rules:

* **Template vs Disk Capacity**:
  Either `template` or `operating_system_disk_capacity` must be provided. If no `template` is given, then `operating_system_disk_capacity` is **required** to manually define the OS disk.

* **Disk Extension**:
  If a template is used, the system provisions a default disk size.
  If you also specify a larger `operating_system_disk_capacity`, Terraform will automatically extend the disk after creation to match your specified size.

* **ISO Support**:
  You can optionally attach an ISO file via `iso_file`.

* **Metadata Application**:
  If `tags`, `description`, or `notes` are provided, these are encoded and submitted as annotate for the VM after creation.

## Argument Reference

### Required

* `client_id` (Number) – ID of the vBridge client owning the VM.
* `name` (String) – The name of the virtual machine.
* `guest_os_id` (String) – vSphere guest OS identifier.
* `cores` (Number) – The number of cores to allocate to the virtual machine.
* `memory_size` (Number) – The amount of memory to allocate to the virtual machine in GB.
* `operating_system_disk_storage_profile` (String) – Disk storage profile (`vStorageT1`, `vStorageT2`, `vStorageT3`).
* `hosting_location_id` (String) – Unique ID for the hosting location.
* `hosting_location_name` (String) – Friendly name for the hosting site.
* `hosting_location_default_network` (String) – Default VM network.
* `backup_type` (String) – Backup strategy (`vBackup`, `vBackupDisk`, or `vBackupNone`).

### Optional

* `template` (String) – Name of a VM template (e.g. `Windows2022_Standard_30GB`).
* `operating_system_disk_capacity` (Number) – OS disk size in GB. Required when no template is provided. Will be used to extend the disk if larger than the template default.
* `iso_file` (String) – Optional ISO file path for provisioning.
* `quote_item` (Map) – Optional cost attribution data.
* `shutdown_protection` (Boolean) – Prevent accidental deletion (`false` by default).
* `tags` (Map of String) – (*Bearer token required*) Tags to attach to the VM.
* `description` (String) – (*Bearer token required*) A freeform description.
* `notes` (String) – (*Bearer token required*) Additional information.

### Lifecycle Configuration

* Use `lifecycle` to ignore known server-side changes such as guest OS ID:

```hcl
lifecycle {
  ignore_changes = [guest_os_id]
}
```

## Attributes Reference

* `vm_id` – Unique ID assigned to the VM by vBridge.
* `mo_ref` – Internal vSphere managed object reference.

## Import

```shell
terraform import vbridge_virtual_machine.example <vm_id>
```

Replace `<vm_id>` with the actual identifier of your VM in vBridge.

