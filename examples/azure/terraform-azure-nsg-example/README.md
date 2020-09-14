# Terraform Azure @Module@ Example

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
1. Make sure [the azure-sdk-for-go versions match](#check-go-dependencies) in [/test/go.mod](/test/go.mod) and in [test/azure/terraform_azure_@module@_test.go](/test/terraform_azure_nic_test.go).
1. `go build terraform_azure_@module@_test.go`
1. `go test -v -run TestTerraformAzure@Module@Example`

## Module test APIs

This module exposes the following top-level API methods:

- GetDefaultNsgRulesClientE returns a rules client that retuns the _default_ NSG rules<br>
`func GetDefaultNsgRulesClientE(subscriptionID string) (network.DefaultSecurityRulesClient, error)`
- GetCustomNsgRulesClientE returns a rules client that returns the _custom_ NSG rules<br>
`func GetCustomNsgRulesClientE(subscriptionID string) (network.SecurityRulesClient, error)`
- GetAllNSGRulesE returns a slice containing all rules from the NSG (including defauls and custom)<br>
`func GetAllNSGRulesE(resourceGroupName, nsgName, subscriptionID string) (NsgRuleSummaryList, error)`

Note that the `GetAllNsgRulesE` method returns an `NsgRuleSummaryList`, which is defined as follows:

- NsgRuleSummaryList holds a colleciton of NsgRuleSummary structs<br>
` type NsgRuleSummaryList struct { SummarizedRules []NsgRuleSummary } `

This collection wrapper has the following method:

- FindRuleByName looks for a matching rule name <br>
`func (summarizedRules *NsgRuleSummaryList) FindRuleByName(t *testing.T, name string) NsgRuleSummary`

Finally, the inner collection is defined as an `NsgRuleSummary` struct, defined as follows:

```
type NsgRuleSummary struct {
	Name                     string
	Description              string
	Protocol                 string
	SourcePortRange          string
	DestinationPortRange     string
	SourceAddressPrefix      string
	DestinationAddressPrefix string
	Access                   string
	Priority                 int32
	Direction                string
}
```

This wrapper has the following methods available:

- AllowsDestinationPort checks to see if the rule allows a specific destination port<br>
`func (summarizedRule *NsgRuleSummary) AllowsDestinationPort(port string) bool`
- AllowsSourcePort checks to see if the rule allows a specific source port<br>
`func (summarizedRule *NsgRuleSummary) AllowsSourcePort(port string) bool`

## Check Go Dependencies

Check that the `github.com/Azure/azure-sdk-for-go` version in your generated `go.mod` for this test matches the version in the terratest [go.mod](https://github.com/gruntwork-io/terratest/blob/master/go.mod) file.  

> This was tested with **go1.14.4**.

### Check Azure-sdk-for-go version

Let's make sure [go.mod](https://github.com/gruntwork-io/terratest/blob/master/go.mod) includes the appropriate [azure-sdk-for-go version](https://github.com/Azure/azure-sdk-for-go/releases/tag/v38.1.0):

```go
require (
    ...
    github.com/Azure/azure-sdk-for-go v38.1.0+incompatible
    ...
)
```

If we make changes to either the **go.mod** or the **go test file**, we should make sure that the go build command works still.

```powershell
go build terraform_azure_@module_test.go
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
