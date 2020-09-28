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

variable "config_name_for_lb01" {
  description = "Frontend Config for Load Balancer 01"
  type        = string
  default     = "terratest-loadbalancer-cfg-01"
}

variable "config_name_for_lb02" {
  description = "Frontend Config for Load Balancer 02"
  type        = string
  default     = "terratest-loadbalancer-cfg-02"
}

variable "loadbalancer01_name" {
  description = "Load balancer with a public IP"
  type        = string
  default     = "terratest-loadbalancer-lb-01"
}

variable "loadbalancer02_name" {
  description = "Load balancer with a private IP"
  type        = string
  default     = "terratest-loadbalancer-lb-02"
}

variable "location" {
  description = "The location where to deploy the resources"
  type        = string
  default     = "East US"
}

variable "pip_for_lb01" {
  description = "Public IP for Load Balancer 01"
  type        = string
  default     = "terratest-loadbalancer-pip-01"
}

variable "privateip_for_lb02" {
  description = "Private IP for Load Balancer 02"
  type        = string
  default     = "10.200.2.10"
}

variable "resource_group_name" {
  description = "The resource group where all resources will be deployed"
  type        = string
  default     = "terratest-loadbalancer-rg"
}

variable "subnet_name" {
  description = "Frontend Subnet for Load Balancer 02"
  type        = string
  default     = "terratest-loadbalancer-snt-02"
}

variable "vnet_name" {
  description = "Virtual Network for Load Balancer 02"
  type        = string
  default     = "terratest-loadbalancer-vnet"
}
