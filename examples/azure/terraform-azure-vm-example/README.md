# Terraform Azure Virtual Machine Example

This folder contains a complete Terraform VM module that deploys resources in [Azure](https://azure.microsoft.com/) to demonstrate
how you can use Terratest to write automated tests for your Azure Virtual Machine Terraform code. This module deploys these resources:

* A [Virtual Machine](https://azure.microsoft.com/en-us/services/virtual-machines/) and gives that VM the following:
    * [Virtual Machine](https://docs.microsoft.com/en-us/azure/virtual-machines/) with the namne specified in the `vm_name` variable.
    * [Managed Disk](https://docs.microsoft.com/en-us/azure/virtual-machines/managed-disks-overview) with the namne specified in the `managed_disk_name` variable.
    * [Availability Set](https://docs.microsoft.com/en-us/azure/virtual-machines/availability) with the namne specified in the `availability_set_name` variable.
* A [Virtual Network](https://azure.microsoft.com/en-us/services/virtual-network/) module that gives the following resources:
    * [Virtual Network](https://docs.microsoft.com/en-us/azure/virtual-network/) with the name specified in the `virtual_network_name` variable.
    * [Subnet](https://docs.microsoft.com/en-us/rest/api/virtualnetwork/subnets) with the name specified in the `subnet_name` variable.
    * [Public Address](https://docs.microsoft.com/en-us/azure/virtual-network/public-ip-addresses) with the name specified in the `public_ip_name` variable.
    * [Network Interface](https://docs.microsoft.com/en-us/azure/virtual-network/virtual-network-network-interface) with the name specified in the `network_interface_name` variable.

Check out [test/azure/terraform_azure_vm_test.go](/test/azure/terraform_azure_vm_test.go) to see how you can write
automated tests for this module.

Note that the Virtual Machine madule creates a Microsoft Windows Server Image with a managed disk, availability set and network configuration for demonstration purposes.

**WARNING**: This module and the automated tests for it deploy real resources into your Azure account which can cost you
money. The resources are all part of the [Azure Free Account](https://azure.microsoft.com/en-us/free/), so if you haven't used that up,
it should be free, but you are completely responsible for all Azure charges.

## Running this module manually

1. Sign up for [Azure](https://azure.microsoft.com/).
1. Configure your Azure credentials using one of the [supported methods for Azure CLI
   tools](https://docs.microsoft.com/en-us/cli/azure/azure-cli-configuration?view=azure-cli-latest).
1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Run `terraform init`.
1. Run `terraform apply`.
1. When you're done, run `terraform destroy`.

## Running automated tests against this module

1. Sign up for [Azure](https://azure.microsoft.com/).
1. Configure your Azure credentials using one of the [supported methods for Azure CLI
   tools](https://docs.microsoft.com/en-us/cli/azure/azure-cli-configuration?view=azure-cli-latest).
1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. [Review environment variables](#review-environment-variables).
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. Make sure [the azure-sdk-for-go versions match](#check-go-dependencies) in [/go.mod](/go.mod) and in [test/azure/terraform_azure_vm_test.go](/test/azure/terraform_azure_vm_test.go).
1. `go build terraform_azure_vm_test.go`
1. `go test -v -run TestTerraformAzureVmExample -timeout 20m` 
    * Note the extra -timeout flag of 20 minutes ensures proper Azure resource removal time.

## Module test APIs

- `VirtualMachineExists` indicates whether the speficied Azure Virtual Machine exists
func VirtualMachineExists(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) bool

- `GetVirtualMachineAdminUser` gets the Admin Username of the specified Azure Virtual Machine
func GetVirtualMachineAdminUser(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string

- `GetVirtualMachineAvailabilitySetID` gets the Availability Set ID of the specified Azure Virtual Machine
func GetVirtualMachineAvailabilitySetID(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string

- `GetVirtualMachineImage` gets the Image of the specified Azure Virtual Machine
func GetVirtualMachineImage(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) VMImage
    ```
    type VMImage struct {
        Publisher string
        Offer     string
        SKU       string
        Version   string
    }
    ```
- `GetVirtualMachineManagedDiskCount` gets the Managed Disk count of the specified Azure Virtual Machine\
func GetVirtualMachineManagedDiskCount(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) int

- `GetVirtualMachineManagedDisks` gets the list of Managed Disk names of the specified Azure Virtual Machine\
func GetVirtualMachineManagedDisks(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) []string

- `GetVirtualMachineNicCount` gets the Network Interface count of the specified Azure Virtual Machine\
func GetVirtualMachineNicCount(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) int

- `GetVirtualMachineNics` gets a list of Network Interface names for a speficied Azure Virtual Machine\
func GetVirtualMachineNics(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) []string

- `GetVirtualMachineOsDiskName` gets the OS Disk name of the specified Azure Virtual Machine\
func GetVirtualMachineOsDiskName(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string

- `GetVirtualMachineSize` gets the Size Type of the specified Azure Virtual Machine\
func GetVirtualMachineSize(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) compute.VirtualMachineSizeTypes

- `GetVirtualMachineState` gets the State of the specified Azure Virtual Machine\
func GetVirtualMachineState(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string

- `GetVirtualMachineTags` gets the Tags of the specified Virtual Machine as a map\
func GetVirtualMachineTags(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) map[string]string

- `GetVirtualMachineInstance` gets a local Virtual Machine instance in the specified Resource Group\
func GetVirtualMachineInstance(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) *Instance

- `GetVirtualMachineInstanceSize` gets the size of the Virtual Machine\
func (vm *Instance) GetVirtualMachineInstanceSize() compute.VirtualMachineSizeTypes

- `GetResourceGroupVirtualMachines` gets a list of all Virtual Machine names in the specified Resource Group\
func GetResourceGroupVirtualMachines(t testing.TestingT, resGroupName string, subscriptionID string) []string

- `GetResourceGroupVirtualMachinesObjects` gets all Virtual Machine objects in the specified Resource Group\
func GetResourceGroupVirtualMachinesObjects(t testing.TestingT, resGroupName string, subscriptionID string) *map[string]compute.VirtualMachineProperties

- `GetVirtualMachine` gets a Virtual Machine in the specified Azure Resource Group\
func GetVirtualMachine(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) *compute.VirtualMachine

- `GetVirtualMachineClientE` creates a Azure Virtual Machine client in the specified Azure Subscription\
func GetVirtualMachineClientE(subscriptionID string) (*compute.VirtualMachinesClient, error)

## Check Go Dependencies

Check that the `github.com/Azure/azure-sdk-for-go` version in your generated `go.mod` for this test matches the version in the terratest [go.mod](https://github.com/gruntwork-io/terratest/blob/master/go.mod) file.  

> This was tested with **go1.14.4**.

### Check Azure-sdk-for-go version

Let's make sure [go.mod](https://github.com/gruntwork-io/terratest/blob/master/go.mod) includes the appropriate [azure-sdk-for-go version](https://github.com/Azure/azure-sdk-for-go/releases/tag/v46.1.0):

```go
require (
    ...
    github.com/Azure/azure-sdk-for-go v46.1.0+incompatible
    ...
)
```

If we make changes to either the **go.mod** or the **go test file**, we should make sure that the go build command works still.

```powershell
go build terraform_azure_vm_test.go
```

## Review Environment Variables

As part of configuring terraform for Azure, we'll want to check that we have set the appropriate [credentials](https://docs.microsoft.com/en-us/azure/terraform/terraform-install-configure?toc=https%3A%2F%2Fdocs.microsoft.com%2Fen-us%2Fazure%2Fterraform%2Ftoc.json&bc=https%3A%2F%2Fdocs.microsoft.com%2Fen-us%2Fazure%2Fbread%2Ftoc.json#set-up-terraform-access-to-azure) and also that we set the [environment variables](https://docs.microsoft.com/en-us/azure/terraform/terraform-install-configure?toc=https%3A%2F%2Fdocs.microsoft.com%2Fen-us%2Fazure%2Fterraform%2Ftoc.json&bc=https%3A%2F%2Fdocs.microsoft.com%2Fen-us%2Fazure%2Fbread%2Ftoc.json#configure-terraform-environment-variables) on the testing host.

```bash
export ARM_CLIENT_ID=your_app_id
export ARM_CLIENT_SECRET=your_password
export ARM_SUBSCRIPTION_ID=your_subscription_id
export ARM_TENANT_ID=your_tenant_id
```

Note, in a Windows environment, these should be set as **system environment variables**.  We can use a PowerShell console with administrative rights to update these environment variables:

```powershell
[System.Environment]::SetEnvironmentVariable("ARM_CLIENT_ID",$your_app_id,[System.EnvironmentVariableTarget]::Machine)
[System.Environment]::SetEnvironmentVariable("ARM_CLIENT_SECRET",$your_password,[System.EnvironmentVariableTarget]::Machine)
[System.Environment]::SetEnvironmentVariable("ARM_SUBSCRIPTION_ID",$your_subscription_id,[System.EnvironmentVariableTarget]::Machine)
[System.Environment]::SetEnvironmentVariable("ARM_TENANT_ID",$your_tenant_id,[System.EnvironmentVariableTarget]::Machine)
```

