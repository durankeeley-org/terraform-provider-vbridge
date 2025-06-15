# Softsource vBridge Terraform Provider

This is a custom Terraform provider for managing infrastructure on [Softsource vBridge](https://www.svbgroup.co.nz/), a New Zealand-based IaaS public cloud provider headquartered in Christchurch.

The provider enables users to deploy and manage virtual machines and related infrastructure through the vBridge API.

---

## Features

* Provision virtual machines on vBridge’s hosting platform
* Manage infrastructure with familiar Terraform syntax
* Supports API key-based authentication
* Easy build/test pipeline via [go-batect](https://github.com/durankeeley/go-batect) (a Go-native replacement for Batect)

---

## Prerequisites

* [Go 1.21+](https://golang.org/dl/)
* [Terraform](https://developer.hashicorp.com/terraform/downloads)
* [`go-batect`](https://github.com/your-org/go-batect) installed or vendored into your project
* API Key + Email from Softsource vBridge

---

## Quick Start

### 1. Clone and Build with `go-batect`

Instead of using Makefiles or shell scripts, use `go-batect` tasks:

```sh
go-batect build
```

This builds provider binaries for both Windows and Linux under the `binaries/` directory.

---

### 2. Install the Provider Locally

After building, copy the binary into Terraform’s plugin directory.

#### Windows

```powershell
Copy-Item binaries\terraform-provider-vbridge.exe `
  "$env:APPDATA\terraform.d\plugins\durankeeley.com\vbridge\vbridge-vm\1.0.1\windows_amd64\terraform-provider-vbridge-vm.exe"
```

#### Linux/macOS

```sh
mkdir -p ~/.terraform.d/plugins/durankeeley.com/vbridge/vbridge-vm/1.0.1/linux_amd64
cp binaries/terraform-provider-vbridge ~/.terraform.d/plugins/durankeeley.com/vbridge/vbridge-vm/1.0.1/linux_amd64/terraform-provider-vbridge-vm
chmod +x ~/.terraform.d/plugins/.../terraform-provider-vbridge-vm
```

---

### 3. Configure & Deploy

1. Copy the secrets example:

   ```sh
   cp secret.tfvars.example secret.tfvars
   ```
2. Initialise Terraform:

   ```sh
   terraform init
   ```
3. Apply:

   ```sh
   terraform apply -var-file="secret.tfvars"
   ```

---

##  Provider Configuration

```hcl
terraform {
  required_providers {
    vbridge = {
      version = "~> 1.0.1"
      source  = "durankeeley.com/vbridge/vbridge-vm"
    }
  }
}

provider "vbridge" {
  auth_type  = var.auth_type
  api_key    = var.api_key
  user_email = var.api_user_email
}
```

Variables (via `secret.tfvars`):

```hcl
auth_type      = "apiKey"
api_key        = "your-api-key"
api_user_email = "your@email.com"
client_id      = 1234
```

---

## Example Resource

```hcl
resource "vbridge_virtual_machine" "example" {
  provider                      = vbridge
  client_id                     = var.client_id
  name                          = "prod-chch"
  template                      = "Windows2022_Standard_30GB"
  guest_os_id                   = "windows2019srv_64Guest"
  cores                         = 2
  memory_size                   = 6
  operating_system_disk_storage_profile = "vStorageT1"
  hosting_location_id           = "vcchcres"
  hosting_location_name         = "Christchurch"
  hosting_location_default_network = "vcchcnet-prod"
  backup_type                   = "vBackupDisk"

  lifecycle {
    ignore_changes = [
      guest_os_id
    ]
  }
}
```

---

## Code Quality & Security

Run all checks via `go-batect`:

```sh
go-batect check-all
```

Tasks included:

* `build` – builds binaries for Windows and Linux
* `test` – runs all Go tests
* `security-scan` – scans the provider source via [Trivy](https://github.com/aquasecurity/trivy)

---

## Debugging Terraform

### Windows PowerShell

```powershell
$env:TF_LOG="DEBUG"
$env:TF_LOG_PATH="C:\temp\terraform.log"
```

### Linux/macOS

```sh
export TF_LOG=DEBUG
export TF_LOG_PATH=/tmp/terraform.log
```

---

## About

This provider is created and maintained by [Duran Keeley](https://github.com/durankeeley) to enable automation on Softsource vBridge’s IaaS platform, with a focus on Go-native workflows and developer-first tooling.

For support or feature requests, feel free to open an issue or contribute via PR.

---

## License

MIT 
