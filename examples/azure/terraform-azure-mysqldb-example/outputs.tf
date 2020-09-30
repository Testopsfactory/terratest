output "resource_group_name" {
  value = azurerm_resource_group.mysql.name
} 

output "mysql_server_name" {
  value = azurerm_mysql_server.mysqlserver.name
}

output "mysql_database_name" {
  value = azurerm_mysql_database.mysqldb.name
}