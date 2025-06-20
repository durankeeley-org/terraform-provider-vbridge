---
page_title: "vbridge_objectstorage_bucket Resource - terraform-provider-vbridge"
subcategory: ""
description: |-
  This resource allows you create a S3 compatible storage bucket in vBridge.
---

# vbridge_virtual_machine (Resource)

This resource allows you to create a S3 compatible storage bucket in vBridge. It can be managed with the AWS CLI or other S3 compatible tools or the AWS terraform provider.

## Example Usage

```terraform
terraform {
  required_providers {
    vbridge = {
      version = "~> 1.0.2"
      source  = "durankeeley-org/vbridge"
    }

    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "vbridge" {
  auth_type  = var.auth_type
  api_key    = var.api_key
  user_email = var.api_user_email
}

variable "objectstorage_tenant_id" {
  description = "Softsource vBridge Object Storage Tenant ID"
  type        = number
  sensitive   = true
}

variable "canonical_user_id" {
  description = "Softsource vBridge Object Storage Canonical User ID"
  type        = string
  sensitive   = true
}

locals {
  subscription = "prod"
  christchurch_hosting = {
    christchurch_shortname = "chch"
    christchurch_location = "vcchcres"
    christchurch_locationname = "Christchurch"
    christchurch_network = "vcchcnet-prod"
    christchurch_s3_endpoint = "https://s3-chc.mycloudspace.co.nz"
  }
}

resource "vbridge_objectstorage_bucket" "examplebucket" {
  objectstorage_tenant_id = var.objectstorage_tenant_id
  canonical_user_id = "${var.canonical_user_id}"
  bucket_name = "sdcterraformbucket"
  object_lock = false
}

provider "aws" {
  access_key                  = ""
  secret_key                  = ""
  region                      = "us-east-1"

  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    s3       = local.christchurch_hosting.christchurch_s3_endpoint
    dynamodb = local.christchurch_hosting.christchurch_s3_endpoint
  }
}

resource "aws_s3_bucket_cors_configuration" "example" {
  bucket = vbridge_objectstorage_bucket.examplebucket.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["${local.christchurch_hosting.christchurch_s3_endpoint}"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }

  cors_rule {
    allowed_methods = ["GET"]
    allowed_origins = ["*"]
  }
}

```