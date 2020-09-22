# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE AVAILABILITY SET
# This is an example of how to deploy an Azure Availability Set with a Virtual Machine in the availability set 
# and the minimum network resources for the VM.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_availabilityset_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

provider "azurerm" {
  version = "=2.20.0"
  features {}
}

# ---------------------------------------------------------------------------------------------------------------------
# PIN TERRAFORM VERSION TO >= 0.12
# The examples have been upgraded to 0.12 syntax
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  # This module is now only being tested with Terraform 0.13.x. However, to make upgrading easier, we are setting
  # 0.12.26 as the minimum version, as that version added support for required_providers with source URLs, making it
  # forwards compatible with 0.13.x code.
  required_version = ">= 0.12.26"
}

# ---------------------------------------------------------------------------------------------------------------------
# GENERATE RANDOMIZATION STRINGS
# These random strings help prevent resource name collision and improve test security
# ---------------------------------------------------------------------------------------------------------------------

resource "random_string" "avsexample" {
  length  = 3
  lower   = true
  upper   = false
  number  = false
  special = false
}

resource "random_password" "avsexample" {
  length           = 16
  override_special = "-_%@"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "avsexample" {
  name     = format("%s-%s-%s", "terratest", random_string.avsexample.result, "rg")
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY THE AVAILABILITY SET
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_availability_set" "avsexample" {
  name                        = format("%s-%s", random_string.avsexample.result, "avs")
  location                    = azurerm_resource_group.avsexample.location
  resource_group_name         = azurerm_resource_group.avsexample.name
  platform_fault_domain_count = var.avs_fault_domain_count
  managed                     = true
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY MINIMAL NETWORK RESOURCES FOR VM
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_network" "avsexample" {
  name                = format("%s-%s", random_string.avsexample.result, "net")
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.avsexample.location
  resource_group_name = azurerm_resource_group.avsexample.name
}

resource "azurerm_subnet" "avsexample" {
  name                 = format("%s-%s", random_string.avsexample.result, "subnet")
  resource_group_name  = azurerm_resource_group.avsexample.name
  virtual_network_name = azurerm_virtual_network.avsexample.name
  address_prefixes     = ["10.0.17.0/24"]
}

resource "azurerm_network_interface" "avsexample" {
  name                = format("%s-%s", random_string.avsexample.result, "nic")
  location            = azurerm_resource_group.avsexample.location
  resource_group_name = azurerm_resource_group.avsexample.name

  ip_configuration {
    name                          = format("%s-%s", random_string.avsexample.result, "config01")
    subnet_id                     = azurerm_subnet.avsexample.id
    private_ip_address_allocation = "Dynamic"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY VIRTUAL MACHINE
# This VM does not actually do anything and is the smallest size VM available with an Ubuntu image
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_machine" "avsexample" {
  name                             = format("%s-%s", random_string.avsexample.result, "vm")
  location                         = azurerm_resource_group.avsexample.location
  resource_group_name              = azurerm_resource_group.avsexample.name
  network_interface_ids            = [azurerm_network_interface.avsexample.id]
  availability_set_id              = azurerm_availability_set.avsexample.id
  vm_size                          = "Standard_B1ls"
  delete_os_disk_on_termination    = true
  delete_data_disks_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = format("%s-%s", random_string.avsexample.result, "osdisk")
    caching           = "None"
    create_option     = "FromImage"
    managed_disk_type = "StandardSSD_LRS"
  }

  os_profile {
    computer_name  = format("%s-%s", random_string.avsexample.result, "vm")
    admin_username = "testadmin"
    admin_password = random_password.avsexample.result
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}
