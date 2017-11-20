# Configure the AWS Provider
provider "aws" {
    region = "${var.aws_region}"
}

# Create a Security Group
resource "aws_security_group" "demo" {
    name = "terraform-test-${var.ec2_instance_name}"
    description = "Demo Security Group"

    # By not specifying a "vpc_id" we use the default VPC for this region

    # Enable all outbound connections
    egress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }

    # Enable SSH inbound
    ingress {
        from_port = 22
        to_port = 22
        protocol = "tcp"
        cidr_blocks = ["3.3.3.3/32"]
    }

    # Enable SSH inbound
    ingress {
        from_port = 80
        to_port = 80
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }
}

# Launch an EC2 Instance
resource "aws_instance" "demo" {
    ami = "${var.ec2_image}"
    instance_type = "t2.micro"
    key_name = "key-that-does-not-exist"
    vpc_security_group_ids = [ "${aws_security_group.demo.id}" ]
    associate_public_ip_address = true
    tags {
        Name = "terraform-test-${var.ec2_instance_name}"
    }
}