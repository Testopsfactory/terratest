# Terraform Azure NSG Example

This folder contains a simple Terraform module that deploys resources in [Azure](https://azure.microsoft.com/) to demonstrate how you can use Terratest to write automated tests for your Azure Terraform code. This module deploys the following:

* A [Virtual Machine](https://azure.microsoft.com/en-us/services/virtual-machines/) that gives the module the following:
    * [Virtual Machine](https://docs.microsoft.com/en-us/azure/virtual-machines/) with the value specified in the `vm_name` variable.
    * A [Network Security Group](https://docs.microsoft.com/en-us/azure/virtual-network/network-security-groups-overview) created with a single custom rule to allow SSH (port 22) with the nsg name specified in the `nsg_name` variable.

Check out [test/azure/terraform_azure_nsg_example_test.go](/test/azure/terraform_azure_nsg_example_test.go) to see how you can write
automated tests for this module.

Note that the resources deployed in this module don't actually do anything; it just runs the resources for
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
1. Make sure [the azure-sdk-for-go versions match](#check-go-dependencies) in [/test/go.mod](/test/go.mod) and in [test/azure/terraform_azure_nsg_example_test.go](/test/terraform_azure_nsg_example_test.go).
1. `go build terraform_azure_nsg_example_test.go`
1. `go test -v -run TestTerraformAzureNsgExample`