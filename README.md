# vBridge Terraform Provider

## Compile Provider on Windows
```
$env:GOOS = "windows"; $env:GOARCH = "amd64"; go build -o "..\provider-compiled\terraform-provider-vbridge-vm.exe"
$env:GOOS = "linux"; $env:GOARCH = "amd64"; go build -o "..\provider-compiled\terraform-provider-vbridge-vm"
```

## Compile Provider on Linux
```
env GOOS=windows GOARCH=amd64 go build -o "../provider-compiled/terraform-provider-vbridge-vm.exe"
env GOOS=linux GOARCH=amd64 go build -o "../provider-compiled/terraform-provider-vbridge-vm"

```

Copy the binary to the following location

## Install Windows

```
%APPDATA%\terraform.d\plugins\durankeeley.com\vbridge\vbridge-vm\1.0.1\windows_amd64\terraform-provider-vbridge-vm.exe
```


## Install Linux

```
~/.terraform.d/plugins/durankeeley.com/vbridge/vbridge-vm/1.0.1/linux_amd64/terraform-provider-vbridge-vm
```

## Deploy Configuration
Copy the ```secret.tfvars.example``` to ```secret.tfvars```
To install the provider and dependancies use ```terraform init``` and then ```terraform apply -var-file="secret.tfvars"```

### Debug Terraform

```
$env:TF_LOG="DEBUG"
$env:TF_LOG_PATH="C:\temp\terraform.log"
```