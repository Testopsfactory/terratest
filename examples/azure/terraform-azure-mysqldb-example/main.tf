# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE MySQL Database
# This is an example of how to deploy an Azure Mysql database.
# ---------------------------------------------------------------------------------------------------------------------


# ------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
# ------------------------------------------------------------------------------

provider "azurerm" {
  version = "~>2.8.0"
  features {}
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# See test/terraform_azure_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "mysql" {
  name     = var.resource_group_name
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE MySQL SERVER
# ---------------------------------------------------------------------------------------------------------------------

resource "random_id" "mysqlserver" {
  byte_length = 8
}

resource "random_password" "password" {
  length = 16
  special = true
  override_special = "_%@"
}

resource "azurerm_mysql_server" "mysqlserver" {
  name                = "${var.mysqlserver_name}-${lower(random_id.mysqlserver.hex)}"
  location            = azurerm_resource_group.mysql.location
  resource_group_name = azurerm_resource_group.mysql.name

  administrator_login          = var.mysqlserver_admin_login
  administrator_login_password = random_password.password.result

  sku_name   = "B_Gen5_2"
  storage_mb = 5120
  version    = "5.7"

  auto_grow_enabled                 = true
  backup_retention_days             = 7
  geo_redundant_backup_enabled      = true
  infrastructure_encryption_enabled = true
  public_network_access_enabled     = false
  ssl_enforcement           = true
  ssl_minimal_tls_version_enforced  = "TLS1_2"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE MySQL DATA BASE
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_mysql_database" "mysqldb" {
  name                = var.mysqldb_name
  resource_group_name = azurerm_resource_group.mysql.name
  server_name         = azurerm_mysql_server.mysqlserver.name
  charset             = "utf8"
  collation           = "utf8_unicode_ci"
}