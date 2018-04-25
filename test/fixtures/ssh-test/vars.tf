# ---------------------------------------------------------------------------------------------------------------------
# ENVIRONMENT VARIABLES
# Define these secrets as environment variables
# ---------------------------------------------------------------------------------------------------------------------

# AWS_ACCESS_KEY_ID
# AWS_SECRET_ACCESS_KEY

# ---------------------------------------------------------------------------------------------------------------------
# MODULE PARAMETERS
# These variables are expected to be passed in by the operator
# ---------------------------------------------------------------------------------------------------------------------

variable "aws_region" {
  description = "The AWS region in which all resources will be created"
}

variable "ami" {
  description = "The ID of the AMI to run on each instance in this example"
  # Ubuntu Server 14.04 LTS (HVM), SSD Volume Type in us-east-1
  default = "ami-fce3c696"
}

variable "keypair_name" {
  description = "The name of the Key Pair that can be used to SSH to each instance in this example"
}

variable "vpc_id" {
  description = "The ID of the VPC in which to run these instances"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEFINE CONSTANTS
# Generally, these values won't need to be changed.
# ---------------------------------------------------------------------------------------------------------------------

variable "name_prefix" {
  description = "The prefix to use for the names of all resources in these templates"
  default = "ssh-test"
}
