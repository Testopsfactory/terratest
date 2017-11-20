# Configure these variables
variable "aws_region" {
    description = "The AWS region in which all resources are deployed."
}

variable "ec2_key_name" {
    description = "The ssh keypair with which we should launch this EC2 instance."
}

variable "ec2_instance_name" {
    description = "The name of our EC2 instance."
}

variable "ec2_image" {
    description = "The AMI ID to be launched."
}