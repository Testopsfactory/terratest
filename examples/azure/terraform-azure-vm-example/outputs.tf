
output "availability_set_name" {
  value = azurerm_availability_set.avs.name
}

output "managed_disk_name" {
  value = azurerm_managed_disk.disk.name
}

output "network_interface_name" {
  value = azurerm_network_interface.nic.name
}

output "os_disk_name" {
  value = azurerm_virtual_machine.vm_example.storage_os_disk[0].name
}

output "public_ip_name" {
  value = azurerm_public_ip.pip.name
}

output "resource_group_name" {
  value = azurerm_resource_group.vm_rg.name
}

output "subnet_name" {
  value = azurerm_subnet.subnet.name
}

output "virtual_network_name" {
  value = azurerm_virtual_network.vnet.name
}

output "vm_name" {
  value = azurerm_virtual_machine.vm_example.name
}
