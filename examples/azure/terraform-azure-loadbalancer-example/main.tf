# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE AVAILABILITY SET
# This is an example of how to deploy an Azure Availability Set with a Virtual Machine in the availability set 
# and the minimum network resources for the VM.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_loadbalancer_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------
provider "azurerm" {
  version = "~>2.29"
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
  name     = "terratest-lb-rg-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY VIRTUAL NETWORK
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_network" "lb" {
  name                = "vnet-${var.postfix}"
  location            = azurerm_resource_group.lb.location
  resource_group_name = azurerm_resource_group.lb.name
  address_space       = ["10.200.0.0/21"]
  dns_servers         = ["10.200.0.5", "10.200.0.6"]

}

resource "azurerm_subnet" "lb" {
  name                 = "subnet-${var.postfix}"
  resource_group_name  = azurerm_resource_group.lb.name
  virtual_network_name = azurerm_virtual_network.lb.name
  address_prefixes     = ["10.200.2.0/25"]
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY LOAD BALANCER WITH PUBLIC IP 
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_public_ip" "lb" {
  name                    = "pip-${var.postfix}"
  resource_group_name     = azurerm_resource_group.lb.name
  location                = azurerm_resource_group.lb.location
  allocation_method       = "Static"
  ip_version              = "IPv4"
  sku                     = "Basic"
  idle_timeout_in_minutes = "4"
}

resource "azurerm_lb" "lb_public" {
  name                = "lb-public-${var.postfix}"
  location            = azurerm_resource_group.lb.location
  resource_group_name = azurerm_resource_group.lb.name
  sku                 = "Basic"

  frontend_ip_configuration {
    name                 = "config-public"
    public_ip_address_id = azurerm_public_ip.lb.id
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY LOAD BALANCER WITH PRIVATE IP 
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_lb" "lb_private" {
  name                = "lb-private-${var.postfix}"
  location            = azurerm_resource_group.lb.location
  resource_group_name = azurerm_resource_group.lb.name
  sku                 = "Basic"

  frontend_ip_configuration {
    name                          = "config-private-01"
    subnet_id                     = azurerm_subnet.lb.id
    private_ip_address            = var.lb_private_ip
    private_ip_address_allocation = "Static"
  }

  frontend_ip_configuration {
    name                          = "config-private-02"
    subnet_id                     = azurerm_subnet.lb.id
    private_ip_address_allocation = "Dynamic"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY LOAD BALANCER WITH NO FRONTEND CONFIGURATION
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_lb" "lb_no_fe_config" {
  name                = "lb-no-frontend-${var.postfix}"
  location            = azurerm_resource_group.lb.location
  resource_group_name = azurerm_resource_group.lb.name
  sku                 = "Basic"
}
