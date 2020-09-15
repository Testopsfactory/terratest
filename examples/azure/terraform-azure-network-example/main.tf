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
# DEPLOY A RESOURCE GROUP
# See test/terraform_azure_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "main" {
  name     = "${var.prefix}-resources"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY VIRTUAL NETWORK
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_network" "main" {
  name                = "${var.prefix}-network"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  address_space       = ["10.0.0.0/16"]
  dns_servers         = [var.dns_ip_01, var.dns_ip_02]
}

resource "azurerm_subnet" "main" {
  name                 = "${var.prefix}-subnet"
  resource_group_name  = azurerm_resource_group.main.name
  virtual_network_name = azurerm_virtual_network.main.name
  address_prefixes       = [var.subnet_prefix]
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY INTERNAL NETWORK INTERFACE
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_network_interface" "internal" {
  name                = "${var.prefix}-nic-01"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name

  ip_configuration {
    name                          = "terratestconfiguration1"
    subnet_id                     = azurerm_subnet.main.id
    private_ip_address_allocation = "Static"
    private_ip_address            = var.private_ip
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY EXTERNAL NETWORK INTERFACE AND PUBLIC IP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_public_ip" "main" {
  name                    = "${var.prefix}-pip"
  resource_group_name     = azurerm_resource_group.main.name
  location                = azurerm_resource_group.main.location
  allocation_method       = "Static"
  ip_version              = "IPv4"
  sku                     = "Standard"
  idle_timeout_in_minutes = "4"
  domain_name_label       = var.domain_name_label
}

resource "azurerm_network_interface" "external" {
  name                = "${var.prefix}-nic-02"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name

  ip_configuration {
    name                          = "terratestconfiguration1"
    subnet_id                     = azurerm_subnet.main.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.main.id
  }
}

