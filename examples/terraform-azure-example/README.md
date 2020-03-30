# Terraform Azure Example

This folder contains a simple Terraform module that deploys resources in [Azure](https://azure.microsoft.com/) to demonstrate
how you can use Terratest to write automated tests for your Azure Terraform code. This module deploys a [Virtual
Machine](https://azure.microsoft.com/en-us/services/virtual-machines/) and gives that VM a `Name` tag with the value specified in the
`vm_name` variable.

Check out [test/terraform_azure_example_test.go](/test/terraform_azure_example_test.go) to see how you can write
automated tests for this module.

Note that the Virtual Machine in this module doesn't actually do anything; it just runs a Vanilla Ubuntu 16.04 image for
demonstration purposes. For slightly more complicated, real-world examples of Terraform modules, see
[terraform-http-example](/examples/terraform-http-example) and [terraform-ssh-example](/examples/terraform-ssh-example).

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
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `go mod init github.com/<yourrepo>/terratest-module`
1. `go build terraform_azure_example_test.go`
1. Make sure [the azure-sdk-for-go versions match](#check-go-dependencies) in the go.mod file and in the test Go file.
1. [Check environment variables](#check-environment-variables).
1. `go test -v -run TestTerraformAzureExample`




## Check Go Dependencies

Check that the /azure-sdk-for-go version must match the version in the terratest mod.go file.  

> This was tested with **go1.14.1**.  We have included a sample **go.mod** to correspond with the terraform_azure_Example_test.go test, but these steps will include details on how to generate the go module and include matching dependencies.

### Creating a new go.mod file

Suppose we create a new **go.mod** file using `go mod init github.com/<yourrepo>/terratest-module` which may generate a file that includes the following dependencies:

```go
module github.com/my-repo/terratest-module

go 1.14

require (
	github.com/Azure/azure-sdk-for-go v40.6.0+incompatible // indirect
	github.com/gruntwork-io/terratest v0.26.1 // indirect
	github.com/stretchr/testify v1.5.1 // indirect
)
```

In this case, the version for the /azure-sdk-for-go version must match the version in the terratest mod.go file.

We should update **go.mod** to use the appropriate **azure-sdk-for-go version**:

```go
module github.com/my-repo/terratest-module

go 1.14

require (
	github.com/Azure/azure-sdk-for-go v38.1.0+incompatible
	github.com/gruntwork-io/terratest v0.26.1 // indirect
	github.com/stretchr/testify v1.5.1 // indirect
)
```

We should check the corresponding **azure-sdk-for-go version** in the import section for the go test.

```go
import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)
```

If we make changes to either the **go.mod** or the **go test file**, we should make sure that the go build command works still:

```powershell
go build terraform_azure_example_test.go
```

## Check Environment Variables

As part of configuring terraform for Azure, we'll want to check that we have the appropriate [credentials](https://docs.microsoft.com/en-us/azure/terraform/terraform-install-configure?toc=https%3A%2F%2Fdocs.microsoft.com%2Fen-us%2Fazure%2Fterraform%2Ftoc.json&bc=https%3A%2F%2Fdocs.microsoft.com%2Fen-us%2Fazure%2Fbread%2Ftoc.json#set-up-terraform-access-to-azure) and also that we set the [environment variables](https://docs.microsoft.com/en-us/azure/terraform/terraform-install-configure?toc=https%3A%2F%2Fdocs.microsoft.com%2Fen-us%2Fazure%2Fterraform%2Ftoc.json&bc=https%3A%2F%2Fdocs.microsoft.com%2Fen-us%2Fazure%2Fbread%2Ftoc.json#configure-terraform-environment-variables) on the testing host.

```bash
export ARM_CLIENT_ID=your_app_id
export ARM_CLIENT_SECRET=your_password
export ARM_SUBSCRIPTION_ID=your_subscription_id
export ARM_TENANT_ID=your_tenant_id
```

Note, in a Windows environment, these should be set as **system environment variables**.  We can use a PowerShell console with administrative rights:

```powershell
[System.Environment]::SetEnvironmentVariable("ARM_CLIENT_ID",$your_app_id,
[System.EnvironmentVariableTarget]::Machine)
[System.Environment]::SetEnvironmentVariable("ARM_CLIENT_SECRET",$your_password,[System.EnvironmentVariableTarget]::Machine)
[System.Environment]::SetEnvironmentVariable("ARM_SUBSCRIPTION_ID",$your_subscription_id,[System.EnvironmentVariableTarget]::Machine)
[System.Environment]::SetEnvironmentVariable("ARM_TENANT_ID",$your_tenant_id,[System.EnvironmentVariableTarget]::Machine)
```

