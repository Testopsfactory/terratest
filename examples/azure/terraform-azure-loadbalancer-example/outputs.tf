output "lb_no_fe_config_name" {
  value = azurerm_lb.lb_no_fe_config.name
}

output "lb_private_name" {
  value = azurerm_lb.lb_private.name
}

output "lb_private_fe_config_name" {
  value = azurerm_lb.lb_private.frontend_ip_configuration[0].name
}

output "lb_private_ip" {
  value = azurerm_lb.lb_private.frontend_ip_configuration[0].private_ip_address
}

output "lb_public_name" {
  value = azurerm_lb.lb_public.name
}

output "lb_public_fe_config_name" {
  value = azurerm_lb.lb_public.frontend_ip_configuration[0].name
}

output "public_address_name" {
  value = azurerm_public_ip.lb.name
}

output "resource_group_name" {
  value = azurerm_resource_group.lb.name
}
