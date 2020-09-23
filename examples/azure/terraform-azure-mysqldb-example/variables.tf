# ---------------------------------------------------------------------------------------------------------------------
# ENVIRONMENT VARIABLES
# Define these secrets as environment variables
# ---------------------------------------------------------------------------------------------------------------------

# ARM_CLIENT_ID
# ARM_CLIENT_SECRET
# ARM_SUBSCRIPTION_ID
# ARM_TENANT_ID

# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable resource_group_name {
  description = "The name of the resource group in which to create the Microsoft MySQL Server."
  type        = string
  default     = "mysqldatabase-rg"
}

variable location {
  description = "The supported azure location where the resource exists"
  type        = string
  default     = "West US2"
}

variable mysqlserver_name {
    description = "The name of the Microsoft MySQL Server."
    type        = string
    default     = "msmysqlserver"
}

variable mysqlserver_admin_login {
    description = "The administrator login name for the mysql server."
    type        = string
    default     = "mysqladminun"
}

variable mysqldb_name {
    description = "The name of the Microsoft mySQL database."
    type        = string
    default     = "msmysqldatabase"

}