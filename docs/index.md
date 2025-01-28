---
page_title: "Softsource vBridge Provider"
subcategory: ""
description: |-
---


# vbridge Provider

The `vbridge` provider is used to interact with the resources supported by Softsource vBridge IaaS.

Use the navigation to the left to read about the available resources.

The below example shows how to configure the vbridge provider, by passing in the required authentication details as variables from a `secret.tfvars` file e.g. `terraform apply -var-file="secret.tfvars"`

```
auth_type = "apiKey"
api_key = "fakekey"
api_user_email = "example@example.com"
```

## Provider Configuration

```terraform
terraform {
  required_providers {
    vbridge = {
      version = "~> 1.0.2"
      source  = "durankeeley-org/vbridge"
    }
  }
}

provider "vbridge" {
  auth_type  = var.auth_type
  api_key    = var.api_key
  user_email = var.api_user_email
}
```


## Example Usage

```terraform
locals {
  subscription = "prod"
  christchurch_hosting = {
    christchurch_shortname = "chch"
    christchurch_location = "vcchcres"
    christchurch_locationname = "Christchurch"
    christchurch_network = "vcchcnet-prod"
  }
}

variable "auth_type" {
  description = "Softsource vBridge Authentication type to use"
  type        = string
  default     = "apiKey"
}

variable "api_key" {
  description = "Softsource vBridge API Key or Bearer Token"
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

# Configure provider
provider "vbridge" {
  auth_type  = var.auth_type
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
