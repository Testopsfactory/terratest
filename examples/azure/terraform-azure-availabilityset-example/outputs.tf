output "resource_group_name" {
  value = azurerm_resource_group.avsexample.name
}

output "availability_set_name" {
  value = azurerm_availability_set.avsexample.name
}

output "availability_set_fdc" {
  value = azurerm_availability_set.avsexample.platform_fault_domain_count
}

output "vm_name" {
  value = azurerm_virtual_machine.avsexample.name
}
