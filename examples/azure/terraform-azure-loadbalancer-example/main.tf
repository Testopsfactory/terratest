provider "azurerm" {
  version = "~>2.20"
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
# DEPLOY A RESOURCE GROUP
# See test/terraform_azure_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "lb" {
  name     =  "${var.resource_group_name}"
  location = "East US"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY LOAD BALANCER WITH PUBLIC IP 
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_public_ip" "lb" {
  name     =  "${var.pip_forlb01}"
  resource_group_name = azurerm_resource_group.lb.name
  location            = azurerm_resource_group.lb.location
  allocation_method       = "Static"
  ip_version              = "IPv4"
  sku                     = "Basic"
  idle_timeout_in_minutes = "4"
}

resource "azurerm_lb" "lb_public" {
  name     =  "${var.loadbalancer01_name}"
  location            = azurerm_resource_group.lb.location
  resource_group_name = azurerm_resource_group.lb.name
  sku                 = "Basic"

    frontend_ip_configuration {
      name     =  "${var.lb01_feconfig}"
      public_ip_address_id = azurerm_public_ip.lb.id
    }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY LOAD BALANCER WITH PRIVATE IP 
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_network" "lb" {
  name     =  "${var.vnet_name}"
  location            = azurerm_resource_group.lb.location
  resource_group_name = azurerm_resource_group.lb.name
  address_space       = ["10.200.0.0/21"]
  dns_servers         = ["10.200.0.5", "10.200.0.6"]

}

resource "azurerm_subnet" "lb" {
  name     =  "${var.feSubnet_forlb02}"
  resource_group_name = azurerm_resource_group.lb.name
  virtual_network_name = azurerm_virtual_network.lb.name
  address_prefix     = "10.200.2.0/25"
}

resource "azurerm_lb" "lb_private" {
  name     =  "${var.loadbalancer02_name}"
  location            = azurerm_resource_group.lb.location
  resource_group_name = azurerm_resource_group.lb.name
  sku                 = "Basic"

    frontend_ip_configuration {
      name     =  "${var.feIPConfig_forlb02}"
      subnet_id                     = azurerm_subnet.lb.id
      private_ip_address            = "10.200.2.10"
      private_ip_address_allocation = "Static"
    }
}