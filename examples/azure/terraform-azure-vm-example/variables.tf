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

variable "disk_type" {
  description = "temp"
  type        = string
  default     = "StandardSSD_LRS"
}

variable "location" {
  description = "The Azure location to deploy resources too"
  type        = string
  default     = "East US"
}

variable "password" {
  description = "the password to configure for ssh access"
  type        = string
  default     = "horriblepassword1234!"
}

variable "prefix" {
  description = "The prefix that will be attached to all resources deployed"
  type        = string
  default     = "terratest-vm"
}

variable "private_ip" {
  description = "The Static Private IP for the Internal NIC"
  type        = string
  default     = "10.0.17.4"
}

variable "subnet_prefix" {
  description = "The subnet range of IPs for the Virtual Network"
  type        = string
  default     = "10.0.17.0/24"
}

variable "user_name" {
  description = "The username to be provisioned into the vm"
  type        = string
  default     = "testadmin"
}

variable "vm_image_sku" {
  description = "The storage image reference SKU from which the VM is created"
  type        = string
  default     = "2016-Datacenter"
}

variable "vm_image_version" {
  description = "The storage image reference Version from which the VM is created"
  type        = string
  default     = "latest"
}

variable "vm_size" {
  description = "The Azure VM Size of the VM"
  type        = string
  default     = "Standard_B1s"
}
