output "vm_name" {
  value = azurerm_virtual_machine.main.name
}

output "resource_group_name" {
  value = azurerm_resource_group.main.name
}

output "nsg_name" {
  value = azurerm_network_security_group.main.name
}

output "ssh_rule_name" {
  value = azurerm_network_security_rule.allowSSH.name
}

output "http_rule_name" {
  value = azurerm_network_security_rule.blockHTTP.name
}
