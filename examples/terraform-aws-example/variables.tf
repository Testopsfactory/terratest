variable "aws_region" {
  description = "The AWS region to deploy into"
  default     = "us-east-1"
}

variable "instance_name" {
  description = "The Name tag to set for the EC2 Instance."
  default     = "terratest-example"
}