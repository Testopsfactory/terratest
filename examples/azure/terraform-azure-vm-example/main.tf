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
# GENERATE RANDOMIZATION
# This helps avoid resource name collisions and improve test security
# ---------------------------------------------------------------------------------------------------------------------

resource "random_password" "vmexample" {
  length           = 16
  override_special = "-_%@"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY RESOURCE GROUP TO CONTAIN TEST RESOURCES
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "vmexample" {
  name     = "${var.prefix}-resources"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY VIRTUAL NETWORK RESOURCES
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_network" "vmexample" {
  name                = "${var.prefix}-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.vmexample.location
  resource_group_name = azurerm_resource_group.vmexample.name
}

resource "azurerm_subnet" "vmexample" {
  name                 = "${var.prefix}-subnet"
  resource_group_name  = azurerm_resource_group.vmexample.name
  virtual_network_name = azurerm_virtual_network.vmexample.name
  address_prefixes     = [var.subnet_prefix]
}

resource "azurerm_public_ip" "vmexample" {
  name                    = "${var.prefix}-pip"
  resource_group_name     = azurerm_resource_group.vmexample.name
  location                = azurerm_resource_group.vmexample.location
  allocation_method       = "Static"
  ip_version              = "IPv4"
  sku                     = "Standard"
  idle_timeout_in_minutes = "4"
}

resource "azurerm_network_interface" "vmexample" {
  name                = "${var.prefix}-nic"
  location            = azurerm_resource_group.vmexample.location
  resource_group_name = azurerm_resource_group.vmexample.name

  ip_configuration {
    name                          = "terratestconfiguration1"
    subnet_id                     = azurerm_subnet.vmexample.id
    private_ip_address_allocation = "Static"
    private_ip_address            = var.private_ip
    public_ip_address_id          = azurerm_public_ip.vmexample.id
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AVAILABILITY SET
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_availability_set" "vmexample" {
  name                             = "${var.prefix}-avs"
  location                         = azurerm_resource_group.vmexample.location
  resource_group_name              = azurerm_resource_group.vmexample.name
  platform_fault_dovmexample_count = 2
  managed                          = true
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A VIRTUAL MACHINE RUNNING WINDOWS SERVER
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_machine" "vmexample" {
  name                             = "${var.prefix}-vm"
  location                         = azurerm_resource_group.vmexample.location
  resource_group_name              = azurerm_resource_group.vmexample.name
  network_interface_ids            = [azurerm_network_interface.vmexample.id]
  availability_set_id              = azurerm_availability_set.vmexample.id
  vm_size                          = var.vm_size
  license_type                     = "Windows_Server"
  delete_os_disk_on_termination    = true
  delete_data_disks_on_termination = true

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = var.vm_image_sku
    version   = var.vm_image_version
  }

  storage_os_disk {
    name              = "${var.prefix}-osdisk"
    caching           = "None"
    create_option     = "FromImage"
    managed_disk_type = var.disk_type
  }

  os_profile {
    computer_name  = "vmexample"
    admin_username = var.user_name
    admin_password = random_password.vmexample.result
  }
  os_profile_windows_config {
    provision_vm_agent = true
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AND ATTACH MANAGED DISK TO VIRTUAL MACHINE
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_managed_disk" "vmexample" {
  name                 = "${var.prefix}-disk"
  location             = azurerm_resource_group.vmexample.location
  resource_group_name  = azurerm_resource_group.vmexample.name
  storage_account_type = var.disk_type
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_virtual_machine_data_disk_attachment" "vmexample" {
  managed_disk_id    = azurerm_managed_disk.vmexample.id
  virtual_machine_id = azurerm_virtual_machine.vmexample.id
  caching            = "ReadWrite"
  lun                = 10
}
