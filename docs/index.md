---
page_title: "Softsource vBridge Provider"
subcategory: ""
description: |-
---


# vbridge Provider


## Example Usage

```terraform
variable "api_key" {
  description = "Softsource vBridge API Key"
  type        = string
  sensitive   = true
}

variable "api_user_email" {
  description = "Softsource vBridge user email address"
  type        = string
  sensitive   = true
}

variable "client_id" {
  description = "Softsource vBridge client id"
  type        = number
  default     = 0000
}

variable "api_url" {
  description = "URL for the API"
  type        = string
  default     = "https://api.mycloudspace.co.nz/"
}

# Configure provider
provider "vbridge" {
  api_url    = var.api_url
  api_key    = var.api_key
  user_email = var.api_user_email
}

# Create a machine
resource "vbridge_virtual_machine" "example" {
  provider    = vbridge
  client_id   = var.client_id
  name        = "${local.subscription}-${local.christchurch_hosting.christchurch_shortname}"
  template    = "Windows2022_Standard_30GB"
  guest_os_id = "windows2019srv_64Guest"
  cores       = 2
  memory_size = 6
  # operating_system_disk_capacity     = 30 
  operating_system_disk_storage_profile = "vStorageT1"
  iso_file                              = ""
  quote_item                            = {}
  hosting_location_id                   = local.christchurch_hosting.christchurch_location
  hosting_location_name                 = local.christchurch_hosting.christchurch_locationname
  hosting_location_default_network      = local.christchurch_hosting.christchurch_network
  backup_type                           = "vBackupDisk"

  lifecycle {
    ignore_changes = [
      guest_os_id
    ]
  }
}
```
