# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE AVAILABILITY SET
# This is an example of how to deploy an Azure Availability Set with a Virtual Machine in the availability set 
# and the minimum network resources for the VM.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_loadbalancer_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------
provider "azurerm" {
  version = "~>2.20"
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
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "lb" {
  name     = var.resource_group_name
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY VIRTUAL NETWORK
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_network" "lb" {
  name                = var.vnet_name
  location            = azurerm_resource_group.lb.location
  resource_group_name = azurerm_resource_group.lb.name
  address_space       = ["10.200.0.0/21"]
  dns_servers         = ["10.200.0.5", "10.200.0.6"]

}

resource "azurerm_subnet" "lb" {
  name                 = var.subnet_name
  resource_group_name  = azurerm_resource_group.lb.name
  virtual_network_name = azurerm_virtual_network.lb.name
  address_prefixes     = ["10.200.2.0/25"]
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY LOAD BALANCER 01 WITH PUBLIC IP 
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_public_ip" "lb" {
  name                    = var.pip_for_lb01
  resource_group_name     = azurerm_resource_group.lb.name
  location                = azurerm_resource_group.lb.location
  allocation_method       = "Static"
  ip_version              = "IPv4"
  sku                     = "Basic"
  idle_timeout_in_minutes = "4"
}

resource "azurerm_lb" "lb_public" {
  name                = var.loadbalancer01_name
  location            = azurerm_resource_group.lb.location
  resource_group_name = azurerm_resource_group.lb.name
  sku                 = "Basic"

  frontend_ip_configuration {
    name                 = var.config_name_for_lb01
    public_ip_address_id = azurerm_public_ip.lb.id
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY LOAD BALANCER 02 WITH PRIVATE IP 
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_lb" "lb_private" {
  name                = var.loadbalancer02_name
  location            = azurerm_resource_group.lb.location
  resource_group_name = azurerm_resource_group.lb.name
  sku                 = "Basic"

  frontend_ip_configuration {
    name                          = var.config_name_for_lb02
    subnet_id                     = azurerm_subnet.lb.id
    private_ip_address            = var.privateip_for_lb02
    private_ip_address_allocation = "Static"
  }
}
