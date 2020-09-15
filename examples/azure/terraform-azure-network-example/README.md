# Terraform Azure Network Example

This folder contains a simple Terraform module that deploys resources in [Azure](https://azure.microsoft.com/) to demonstrate
how you can use Terratest to write automated tests for your Azure Terraform code. This module deploys to a Virtual Network two Network Interface Cards, one with an internal only IP and another with an internal and external Public IP.

* A [Virtual Network](https://azure.microsoft.com/en-us/services/virtual-network/) module that gives the following resources:
    * [Virtual Network](https://docs.microsoft.com/en-us/azure/virtual-network/) with the name specified in the `virtual_network_name` variable.
    * [Subnet](https://docs.microsoft.com/en-us/rest/api/virtualnetwork/subnets) with the name specified in the `subnet_name` variable.
    * [Public Address](https://docs.microsoft.com/en-us/azure/virtual-network/public-ip-addresses) with the name specified in the `public_ip_name` variable.
    * [Internal Network Interface](https://docs.microsoft.com/en-us/azure/virtual-network/virtual-network-network-interface) with the name specified in the `network_interface_internal` variable.
    * [ExternalNetwork Interface](https://docs.microsoft.com/en-us/azure/virtual-network/virtual-network-network-interface) with the name specified in the `network_interface_external` variable.

Check out [test/azure/terraform_azure_network_test.go](/test/azure/terraform_azure_network_example_test.go) to see how you can write
automated tests for this module.

Note that the Azure Virtual Network, Subnet, Network Interface and Public IP resources in this module don't actually do anything; it just runs the resources for
demonstration purposes.

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
1. Make sure [the azure-sdk-for-go versions match](#check-go-dependencies) in [/go.mod](/go.mod) and in [test/terraform_azure_network_example_test.go](/test/azure/terraform_azure_network_example_test.go).
1. `go build terraform_azure_network_example_test.go`
1. `go test -v -run TestTerraformAzureNetworkExample`

## Module test APIs

These modules expose the following API methods:

### Virtual Networks and Subnets
- `VirtualNetworkExists` indicates whether the speficied Azure Virtual Network exists\
func VirtualNetworkExists(t testing.TestingT, vnetName string, resGroupName string, subscriptionID string) bool

- `SubnetExists` indicates whether the speficied Azure Virtual Network Subnet exists\
func SubnetExists(t testing.TestingT, subnetName string, vnetName string, resGroupName string, subscriptionID string) bool

- `CheckSubnetContainsIP` checks if the Private IP is contined in the Subnet Address Range\
func CheckSubnetContainsIP(t testing.TestingT, IP string, subnetName string, vnetName string, resGroupName string, subscriptionID string) bool

- `GetSubnetIPRange` gets the IPv4 Range of the specified Subnet\
func GetSubnetIPRange(t testing.TestingT, subnetName string, vnetName string, resGroupName string, subscriptionID string) string

- `GetVirtualNetworkDNSServerIPs` gets a list of all Virtual Network DNS server IPs\
func GetVirtualNetworkDNSServerIPs(t testing.TestingT, vnetName string, resGroupName string, subscriptionID string) []string

- `GetVirtualNetworkSubnets` gets all Subnet names and their respective addres prefixes in the specified Virtual Network\
func GetVirtualNetworkSubnets(t testing.TestingT, vnetName string, resGroupName string, subscriptionID string) map[string]string

- `GetVirtualNetworkE` gets Virtual Network in the specified Azure Resource Group\
func GetVirtualNetworkE(t testing.TestingT, vnetName string, resGroupName string, subscriptionID string) (*network.VirtualNetwork, error)

- `GetVirtualNetworksClientE` creates a virtual network client in the specified Azure Subscription\
func GetVirtualNetworksClientE(subscriptionID string) (*network.VirtualNetworksClient, error)

- `GetSubnetE` gets a subnet\
func GetSubnetE(t testing.TestingT, subnetName string, vnetName string, resGroupName string, subscriptionID string) (*network.Subnet, error)

- `GetSubnetClientE` creates a subnet client\
func GetSubnetClientE(subscriptionID string) (*network.SubnetsClient, error)



### Network Interfaces
-    `NetworkInterfaceExists` indicates whether the speficied Azure Network Interface exists\
func NetworkInterfaceExists(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) bool

-    `GetNetworkInterfacePrivateIPs` gets a list of the Private IPs of a Network Interface configs\
func GetNetworkInterfacePrivateIPs(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) []string

-    `GetNetworkInterfacePublicIPs` returns a list of all the Public IPs found in the Network Interface configurations\
func GetNetworkInterfacePublicIPs(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) []string

-    `GetNetworkInterfaceE` gets a Network Interface in the specified Azure Resource Group\
func GetNetworkInterfaceE(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) (*network.Interface, error)

-    `GetNetworkInterfaceClientE` creates a new Network Interface client in the specified Azure Subscription\
func GetNetworkInterfaceClientE(subscriptionID string) (*network.InterfacesClient, error)

### Public Addresses
-    `PublicAddressExists` indicates whether the speficied AzurePublic Address exists\
func PublicAddressExists(t testing.TestingT, publicAddressName string, resGroupName string, subscriptionID string) bool

-    `GetPublicAddressIP` gets the IP of a Public IP Address\
func GetPublicAddressIP(t testing.TestingT, publicAddressName string, resGroupName string, subscriptionID string) string

-    `CheckPublicDNSNameAvailability` checks whether a Domain Name in the cloudapp.azure.com zone is available for use\
func CheckPublicDNSNameAvailability(t testing.TestingT, location string, domainNameLabel string, subscriptionID string) bool

-    `GetPublicIPAddressE` gets a Public IP Addresses in the specified Azure Resource Group\
func GetPublicIPAddressE(t testing.TestingT, publicIPAddressName string, resGroupName string, subscriptionID string) (*network.PublicIPAddress, error)

-    `GetPublicIPAddressClientE` creates a Public IP Addresses client in the specified Azure Subscription\
func GetPublicIPAddressClientE(subscriptionID string) (*network.PublicIPAddressesClient, error)



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
go build terraform_azure_network_example_test.go
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