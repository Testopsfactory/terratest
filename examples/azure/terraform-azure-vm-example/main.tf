# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_vm_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

provider "azurerm" {
  version = "=2.20.0"
  features {}
}

# ---------------------------------------------------------------------------------------------------------------------
# PIN TERRAFORM VERSION TO >= 0.12
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  required_version = ">= 0.12"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "vm" {
  name     = "terratest-vs-rg-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY NETWORK RESOURCES
# This network includes a public address for integration tests
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_network" "vm" {
  name                = "vnet-${var.postfix}"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.vm.location
  resource_group_name = azurerm_resource_group.vm.name
}

resource "azurerm_subnet" "vm" {
  name                 = "subnet-${var.postfix}"
  resource_group_name  = azurerm_resource_group.vm.name
  virtual_network_name = azurerm_virtual_network.vm.name
  address_prefixes     = [var.subnet_prefix]
}

resource "azurerm_public_ip" "vm" {
  name                    = "pip-${var.postfix}"
  resource_group_name     = azurerm_resource_group.vm.name
  location                = azurerm_resource_group.vm.location
  allocation_method       = "Static"
  ip_version              = "IPv4"
  sku                     = "Standard"
  idle_timeout_in_minutes = "4"
}

resource "azurerm_network_interface" "vm" {
  name                = "nic-${var.postfix}"
  location            = azurerm_resource_group.vm.location
  resource_group_name = azurerm_resource_group.vm.name

  ip_configuration {
    name                          = "terratestconfiguration1"
    subnet_id                     = azurerm_subnet.vm.id
    private_ip_address_allocation = "Static"
    private_ip_address            = var.private_ip
    public_ip_address_id          = azurerm_public_ip.vm.id
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AVAILABILITY SET
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_availability_set" "vm" {
  name                        = "avs-${var.postfix}"
  location                    = azurerm_resource_group.vm.location
  resource_group_name         = azurerm_resource_group.vm.name
  platform_fault_domain_count = 2
  managed                     = true
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY VIRTUAL MACHINE
# This VM does not actually do anything and is the smallest size VM available with a Windows image
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_machine" "vm" {
  name                             = "$vm-${var.postfix}"
  location                         = azurerm_resource_group.vm.location
  resource_group_name              = azurerm_resource_group.vm.name
  network_interface_ids            = [azurerm_network_interface.vm.id]
  availability_set_id              = azurerm_availability_set.vm.id
  vm_size                          = var.vm_size
  license_type                     = var.vm_license_type
  delete_os_disk_on_termination    = true
  delete_data_disks_on_termination = true

  storage_image_reference {
    publisher = var.vm_image_publisher
    offer     = var.vm_image_offer
    sku       = var.vm_image_sku
    version   = var.vm_image_version
  }

  storage_os_disk {
    name              = "osdisk-${var.postfix}"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = var.disk_type
  }

  os_profile {
    computer_name  = "vm-${var.postfix}"
    admin_username = var.user_name
    admin_password = random_password.vm.result
  }
  os_profile_windows_config {
    provision_vm_agent = true
  }

  depends_on = [random_password.vm]
}

# Random password used to improve test security
resource "random_password" "vm" {
  length           = 16
  override_special = "-_%@"
  min_upper        = "1"
  min_lower        = "1"
  min_numeric      = "1"
  min_special      = "1"
}

# ---------------------------------------------------------------------------------------------------------------------
# ATTACH A MANAGED DISK TO THE VIRTUAL MACHINE
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_managed_disk" "vm" {
  name                 = "${var.postfix}-disk"
  location             = azurerm_resource_group.vm.location
  resource_group_name  = azurerm_resource_group.vm.name
  storage_account_type = var.disk_type
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_virtual_machine_data_disk_attachment" "vm" {
  managed_disk_id    = azurerm_managed_disk.vm.id
  virtual_machine_id = azurerm_virtual_machine.vm.id
  caching            = "ReadWrite"
  lun                = 10
}
