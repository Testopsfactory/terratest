output "resource_group_name" {
  value = azurerm_resource_group.lb.name
}

output "loadbalancer01_name" {
  value = azurerm_lb.lb_public.name
}

output "lb01_feconfig" {
  value = azurerm_lb.lb_public.frontend_ip_configuration[0].name
}

output "pip_forlb01" {
  value = azurerm_public_ip.lb.name
}

output "loadbalancer02_name" {
  value = azurerm_lb.lb_private.name
}
output "privateip_forlb02" {
  value = azurerm_lb.lb_private.frontend_ip_configuration[0].private_ip_address
}

output "feSubnet_forlb02" {
  value = azurerm_lb.lb_private.frontend_ip_configuration[0].subnet_id
}
