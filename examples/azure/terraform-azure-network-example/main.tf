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
# See test/azure/terraform_azure_network_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "net" {
  name     = "${var.prefix}-resources"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY VIRTUAL NETWORK
# Note this network dosen't actually do anything and is only created for the example.
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_network" "net" {
  name                = "${var.prefix}-vnet"
  location            = azurerm_resource_group.net.location
  resource_group_name = azurerm_resource_group.net.name
  address_space       = ["10.0.0.0/16"]
  dns_servers         = [var.dns_ip_01, var.dns_ip_02]
}

resource "azurerm_subnet" "net" {
  name                 = "${var.prefix}-subnet"
  resource_group_name  = azurerm_resource_group.net.name
  virtual_network_name = azurerm_virtual_network.net.name
  address_prefixes     = [var.subnet_prefix]
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY INTERNAL NETWORK INTERFACE
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_network_interface" "net01" {
  name                = "${var.prefix}-nic-01"
  location            = azurerm_resource_group.net.location
  resource_group_name = azurerm_resource_group.net.name

  ip_configuration {
    name                          = "terratestconfiguration1"
    subnet_id                     = azurerm_subnet.net.id
    private_ip_address_allocation = "Static"
    private_ip_address            = var.private_ip
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY PUBLIC IP AND EXTERNAL NETWORK INTERFACE
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_public_ip" "net" {
  name                    = "${var.prefix}-pip"
  resource_group_name     = azurerm_resource_group.net.name
  location                = azurerm_resource_group.net.location
  allocation_method       = "Static"
  ip_version              = "IPv4"
  sku                     = "Standard"
  idle_timeout_in_minutes = "4"
  domain_name_label       = var.domain_name_label
}

resource "azurerm_network_interface" "net02" {
  name                = "${var.prefix}-nic-02"
  location            = azurerm_resource_group.net.location
  resource_group_name = azurerm_resource_group.net.name

  ip_configuration {
    name                          = "terratestconfiguration1"
    subnet_id                     = azurerm_subnet.net.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.net.id
  }
}

