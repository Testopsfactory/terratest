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
1. `go mod init terratest-module`
1. `go build terraform_azure_example_test.go`
1. Make sure to [check the depedencies](#check-go-dependencies) match in the go.mod file and in the test Go file.
1. `go test -v -run TestTerraformAzureExample`




## Check Go Dependencies
This was tested with **go1.14.1**.  We have included a sample **go.mod** to correspond with the go test, but these steps will include details on how to generate the file and dependencies.

Suppose we create a new **go.mod** file using `go mod init terratest-module` may generate a file that includes the following dependencies:

```go
module terratest-module

go 1.14

require (
	github.com/Azure/azure-sdk-for-go v40.6.0+incompatible // indirect
	github.com/gruntwork-io/terratest v0.26.1 // indirect
	github.com/stretchr/testify v1.5.1 // indirect
)
```

In this case, the dependency for the **azure-sdk-for-go** (`github.com/Azure/azure-sdk-for-go v40.6.0+incompatible // indirect`) needs to match what's used in the Go source test.

We should update go.mod to use the appropriate azure-sdk-for-go:

```go
module terratest-module

go 1.14

require (
	github.com/Azure/azure-sdk-for-go v38.1.0+incompatible
	github.com/gruntwork-io/terratest v0.26.1 // indirect
	github.com/stretchr/testify v1.5.1 // indirect
)
```

We should check that **azure-sdk-for-go dependency** in the import section for the go test.
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