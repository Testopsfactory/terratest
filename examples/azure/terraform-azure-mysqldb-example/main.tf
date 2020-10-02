# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE MySQL Database
# This is an example of how to deploy an Azure Mysql database.
# ---------------------------------------------------------------------------------------------------------------------


# ------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
# ------------------------------------------------------------------------------

provider "azurerm" {
  version = "~>2.29.0"
  features {}
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# See test/terraform_azure_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "mysql" {
  name     = "terratest-mysql-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE MySQL SERVER
# ---------------------------------------------------------------------------------------------------------------------

resource "random_password" "password" {
  length = 16
  special = true
  override_special = "_%@"
}

resource "azurerm_mysql_server" "mysqlserver" {
  name                = "mysqlserver-${var.postfix}"
  location            = azurerm_resource_group.mysql.location
  resource_group_name = azurerm_resource_group.mysql.name

  administrator_login          = var.mysqlserver_admin_login
  administrator_login_password = random_password.password.result

  sku_name   = var.mysqlserver_sku_name
  storage_mb = var.mysqlserver_storage_mb
  version    = "5.7"

  auto_grow_enabled                 = true
  geo_redundant_backup_enabled      = false
  infrastructure_encryption_enabled = true
  backup_retention_days             = 7
  public_network_access_enabled     = false
  ssl_enforcement_enabled           = true
  ssl_minimal_tls_version_enforced  = "TLS1_2"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE MySQL DATA BASE
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_mysql_database" "mysqldb" {
  name                = "mysqldb-${var.postfix}"
  resource_group_name = azurerm_resource_group.mysql.name
  server_name         = azurerm_mysql_server.mysqlserver.name
  charset             = var.mysqldb_charset
  collation           = var.mysqldb_collation
}
