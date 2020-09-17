provider "azurerm" {
  version = "=2.20.0"
  features {}
}

# ---------------------------------------------------------------------------------------------------------------------
# PIN TERRAFORM VERSION TO >= 0.12
# The examples have been upgraded to 0.12 syntax
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  required_version = ">= 0.12"
}

# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_availabilityset_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------
# GENERATE RANDOMIZATION STRINGS
# This helps avoid resource name collisions and improve test security
# ---------------------------------------------------------------------------------------------------------------------

resource "random_string" "main" {
  length  = 3
  lower   = true
  upper   = false
  number  = false
  special = false
}

resource "random_password" "main" {
  length           = 16
  override_special = "-_%@"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY RESOURCE GROUP TO CONTAIN TEST RESOURCES
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "main" {
  name     = format("%s-%s-%s", "terratest", random_string.main.result, "avs")
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AVAILABILITY SET
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_availability_set" "main" {
  name                        = format("%s-%s", random_string.main.result, "avs")
  location                    = azurerm_resource_group.main.location
  resource_group_name         = azurerm_resource_group.main.name
  platform_fault_domain_count = var.avs_fault_domain_count
  managed                     = true
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY VIRTUAL NETWORK RESOURCES
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_network" "main" {
  name                = format("%s-%s", random_string.main.result, "net")
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
}

resource "azurerm_subnet" "main" {
  name                 = format("%s-%s", random_string.main.result, "subnet")
  resource_group_name  = azurerm_resource_group.main.name
  virtual_network_name = azurerm_virtual_network.main.name
  address_prefixes     = ["10.0.17.0/24"]
}

resource "azurerm_network_interface" "main" {
  name                = format("%s-%s", random_string.main.result, "nic")
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name

  ip_configuration {
    name                          = format("%s-%s", random_string.main.result, "config01")
    subnet_id                     = azurerm_subnet.main.id
    private_ip_address_allocation = "Dynamic"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY VIRTUAL MACHINE TO THE AVAILABILITY SET
# This VM does not actually do anything and is the smallest size VM available with an Ubuntu image
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_machine" "main" {
  name                             = format("%s-%s", random_string.main.result, "vm")
  location                         = azurerm_resource_group.main.location
  resource_group_name              = azurerm_resource_group.main.name
  network_interface_ids            = [azurerm_network_interface.main.id]
  availability_set_id              = azurerm_availability_set.main.id
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
    name              = format("%s-%s", random_string.main.result, "osdisk")
    caching           = "None"
    create_option     = "FromImage"
    managed_disk_type = "StandardSSD_LRS"
  }

  os_profile {
    computer_name  = format("%s-%s", random_string.main.result, "vm")
    admin_username = "testadmin"
    admin_password = random_password.main.result
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}
