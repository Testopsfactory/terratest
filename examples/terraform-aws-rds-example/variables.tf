variable "name" {
  description = "Name of the database"
  default     = "terratest-example"
}

variable "engine_name" {
  description = "Name of the database engine"
  default     = "mysql"
}

variable "major_engine_version" {
  description = "MAJOR.MINOR version of the DB engine"
  default     = "5.7"
}

variable "family" {
  description = "Family of the database"
  default     = "mysql5.7"
}

variable "username" {
  description = "Master username of the DB"
  default     = "username"
}

# Not a good idea to keep in plain text. Would be better to pass through encrypted
# environment variables or through some secret management solution.
variable "password" {
  description = "Master password of the DB"
  default     = "password"
}

variable "allocated_storage" {
  default     = 5
  description = "Disk space to be allocated to the DB instance"
}

variable "license_model" {
  default     = "general-public-license"
  description = "License model of the DB instance"
}

variable "engine_version" {
  default     = "5.7.21"
  description = "Version of the database to be launched"
}
