output "loadbalancer_01_name" {
  value = azurerm_lb.lb_public.name
}

output "loadbalancer_02_name" {
  value = azurerm_lb.lb_private.name
}

output "config_name_for_lb01" {
  value = azurerm_lb.lb_public.frontend_ip_configuration[0].name
}

output "config_name_for_lb02" {
  value = azurerm_lb.lb_private.frontend_ip_configuration[0].name
}

output "public_address_name_for_lb01" {
  value = azurerm_public_ip.lb.name
}

output "privateip_for_lb02" {
  value = azurerm_lb.lb_private.frontend_ip_configuration[0].private_ip_address
}

output "resource_group_name" {
  value = azurerm_resource_group.lb.name
}
