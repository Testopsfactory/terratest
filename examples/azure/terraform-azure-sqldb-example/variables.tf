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
  description = "The name of the resource group in which to create the Microsoft SQL Server."
  type        = string
  default     = "sqldatabase-rg"
}

variable location {
  description = "The supported azure location where the resource exists"
  type        = string
  default     = "West US2"
}

variable sqlserver_name {
    description = "The name of the Microsoft SQL Server."
    type        = string
    default     = "mssqlserver"
}

variable sqlserver_admin_login {
    description = "The administrator login name for the sql server."
    type        = string
    default     = "3cx1n1z5r079b"
}

variable tags {
    description = "A mapping of tags to assign to the resource."
    type        = string
    default = "Development"
}

variable sqldb_name {
    description = "The name of the Microsoft SQL database."
    type        = string
    default     = "mssqldatabase"

}