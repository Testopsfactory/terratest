# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE SQL Database
# This is an example of how to deploy an Azure sql database.
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

resource "azurerm_resource_group" "sql" {
  name     = var.resource_group_name
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE SQL SERVER
# ---------------------------------------------------------------------------------------------------------------------

resource "random_id" "sqlserver" {
  byte_length = 8
}

resource "random_password" "password" {
  length = 16
  special = true
  override_special = "_%@"
}

resource "azurerm_sql_server" "sqlserver" {
  name                         = "${var.sqlserver_name}-${lower(random_id.sqlserver.hex)}"
  resource_group_name          = azurerm_resource_group.sql.name
  location                     = azurerm_resource_group.sql.location
  version                      = "12.0"
  administrator_login          = var.sqlserver_admin_login
  administrator_login_password = random_password.password.result

  tags = {
    environment = var.tags
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE STORAGE ACCOUNT
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_storage_account" "sql" {
  name                     = var.sa_name
  resource_group_name      = azurerm_resource_group.sql.name
  location                 = azurerm_resource_group.sql.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE SQL DATA BASE
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_sql_database" "sqldb" {
  name                = var.sqldb_name
  resource_group_name = azurerm_resource_group.sql.name
  location            = azurerm_resource_group.sql.location
  server_name         = azurerm_sql_server.sqlserver.name

  extended_auditing_policy {
    storage_endpoint                        = azurerm_storage_account.sql.primary_blob_endpoint
    storage_account_access_key              = azurerm_storage_account.sql.primary_access_key
    storage_account_access_key_is_secondary = true
    retention_in_days                       = 6
  }

  tags = {
    environment = var.tags
  }
}