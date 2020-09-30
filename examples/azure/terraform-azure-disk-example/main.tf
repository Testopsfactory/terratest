# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE MANAGED DISK
# This is an example of how to deploy a managed disk.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_disk_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

provider "azurerm" {
  version = "~> 2.29"
  features {}
}

# ---------------------------------------------------------------------------------------------------------------------
# PIN TERRAFORM VERSION TO >= 0.12
# The examples have been upgraded to 0.12 syntax
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  required_version = ">= 0.12.26"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "disk" {
  name     = "terratest-disk-rg-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY THE DISK
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_managed_disk" "disk" {
  name                 = "disk-${var.postfix}"
  location             = azurerm_resource_group.disk.location
  resource_group_name  = azurerm_resource_group.disk.name
  storage_account_type = var.disk_type
  create_option        = "Empty"
  disk_size_gb         = 10
}
