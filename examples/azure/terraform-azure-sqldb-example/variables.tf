# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable resource_group_name {
  description = "The name of the resource group in which to create the Microsoft SQL Server."
  default     = "sqldatabase-rg"
}

variable location {
  description = "The supported azure location where the resource exists"
  default     = "West US2"
}

variable sqlserver_name {
    description = "The name of the Microsoft SQL Server."
    default     = "mssqlserver"
}

variable sqlserver_admin_login {
    description = "The administrator login name for the sql server."
    default     = "3cx1n1z5r079b"
}

variable tags {
    description = "A mapping of tags to assign to the resource."
    default = "Development"
}

variable sa_name {
    description = "The name of the storage account."
    default = "examplesqlsa"
}

variable sqldb_name {
    description = "The name of the Microsoft SQL database."
    default     = "mssqldatabase"

}