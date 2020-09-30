
output "availability_set_name" {
  value = azurerm_availability_set.vm.name
}

output "managed_disk_name" {
  value = azurerm_managed_disk.vm.name
}

output "network_interface_name" {
  value = azurerm_network_interface.vm.name
}

output "os_disk_name" {
  value = azurerm_virtual_machine.vm.storage_os_disk[0].name
}

output "public_ip_name" {
  value = azurerm_public_ip.vm.name
}

output "resource_group_name" {
  value = azurerm_resource_group.vm.name
}

output "subnet_name" {
  value = azurerm_subnet.vm.name
}

output "virtual_network_name" {
  value = azurerm_virtual_network.vm.name
}

output "vm_name" {
  value = azurerm_virtual_machine.vm.name
}
