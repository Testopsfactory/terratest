# Terraform Azure SQL DB Example

This folder contains a Terraform module that deploys a sql database in [Azure](https://azure.microsoft.com/) to demonstrate how you can use Terratest to write automated tests for your Azure Terraform code. 

This module deploys [Azure SQL Database](https://azure.microsoft.com/en-us/services/sql-database/.

Check out [test/azure/terraform_azure_sqldb_example_test.go](./../../../test/azure/terraform_azure_sqldb_example_test.go) to see how you can write automated tests for this module and validate the configuration of the parameters and options. 

**WARNING**: This module and the automated tests for it deploy real resources into your Azure account which can cost you money. 

## Prerequisite: Setup Azure CLI access
1. Sign up for [Azure](https://azure.microsoft.com/).
1. Install [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest)
1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Login to Azure on the CLI with `az login` or `az login --use-device`, and then configure the CLI.

## Running this module manually
1. Create [Service Principal](https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli?view=azure-cli-latest) then set the value to the environment variables. 
1. Run `terraform init`
1. Run `terraform apply`
1. check the services are getting deployed to your azure subscription.
1. When you're done, run `terraform destroy`.

### Example

```bash
$ az login 
$ export ARM_SUBSCRIPTION_ID={YOUR_SUBSCRIPTION_ID} 
$ az ad sp create-for-rbac
$ export TF_VAR_client_id={YOUR_SERVICE_PRINCIPAL_APP_ID}
$ export TF_VAR_client_secret={YOUR_SERVICE_PRINCIPAL_PASSWORD}
$ terraform init
$ terraform apply
$ terraform destroy
```

## Running automated tests against this module
1. Create [Service Principal](https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli?view=azure-cli-latest) then set the value to the environment variables. 
1. Install [Golang](https://golang.org/) version `1.13+` required. 
1. `cd test`
1. `go test -v -timeout 60m -tags azure -run TestTerraformAzureSQLDBExample`
### Example

```bash
$ az login 
$ export ARM_SUBSCRIPTION_ID={YOUR_SUBSCRIPTION_ID} 
$ export TF_VAR_client_id={YOUR_SERVICE_PRINCIPAL_APP_ID}
$ export TF_VAR_client_secret={YOUR_SERVICE_PRINCIPAL_PASSWORD}
$ cd test
$ go test -v -timeout 60m -tags azure -run TestTerraformAzureSQLDBExample
```
