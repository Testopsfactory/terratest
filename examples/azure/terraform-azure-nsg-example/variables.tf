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

variable "resource_group_name" {
  description = "Name for the resource group holding resources for this example"
  type        = string
  default     = "terratest-nsg-example"
}

variable "location" {
  description = "The Azure region in which to deploy this sample"
  type        = string
  default     = "East US"
}

variable "vnet_name" {
  description = "Name for the example virtual network"
  type        = string
  default     = "terratest-nsg-example-network"
}

variable "subnet_name" {
  description = "Name for the example virtual network default subnet"
  type        = string
  default     = "terratest-nsg-example-subnet"
}

variable "vm_nic_name" {
  description = "Name for the NIC attached to our example VM."
  type        = string
  default     = "terratest-nsg-example-nic"
}

variable "vm_nic_ip_config_name" {
  description = "Name for the NIC IP configuration attached to our example VM."
  type        = string
  default     = "terratest-nsg-example-nic-ip-config"
}

variable "nsg_name" {
  description = "Name for the example NSG."
  type        = string
  default     = "terratest-nsg-example-nsg"
}

variable "nsg_rule_name" {
  description = "Name for the example NSG rule used in this example."
  type        = string
  default     = "terratest-nsg-example-nsg-rule"
}

variable "vm_name" {
  description = "The name of the VM used in this example"
  type        = string
  default     = "terratest-nsg-example-vm"
}

variable "vm_size" {
  description = "The size of the VM to deploy"
  type        = string
  default     = "Standard_B1s"
}

variable "hostname" {
  description = "The hostname of the new VM to be configured"
  type        = string
  default     = "terratest-nsg-example-vm"
}

variable "os_disk_name" {
  description = "The of the OS disk to use on our example VM."
  type        = string
  default     = "terratest-nsg-example-os-disk"
}

variable "password" {
  description = "The password to configure for SSH access"
  type        = string
  default     = "HorriblePassword1234!"
}

variable "username" {
  description = "The username to be provisioned into your VM"
  type        = string
  default     = "testadmin"
}
