# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# CREATE TWO EC2 INSTANCES FOR TESTING SSH CONNECTIVITY
# These templates deploy two EC2 instances, one with a public IP and one with only a private IP. These can be used to
# test SSH connectivity.
# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

# ---------------------------------------------------------------------------------------------------------------------
# CONFIGURE OUR AWS CONNECTION
# ---------------------------------------------------------------------------------------------------------------------

provider "aws" {
  region = "${var.aws_region}"
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE AN AWS INSTANCE WITH A PUBLIC IP
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_instance" "example_public" {
  ami = "${var.ami}"
  instance_type = "t2.nano"
  key_name = "${var.keypair_name}"
  vpc_security_group_ids = ["${aws_security_group.example.id}"]
  associate_public_ip_address = true
  tags {
    Name = "${var.name_prefix}-public"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE AN AWS INSTANCE WITH ONLY A PRIVATE IP
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_instance" "example_private" {
  ami = "${var.ami}"
  instance_type = "t2.nano"
  key_name = "${var.keypair_name}"
  vpc_security_group_ids = ["${aws_security_group.example.id}"]
  associate_public_ip_address = false
  tags {
    Name = "${var.name_prefix}-private"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE A SECURITY GROUP THAT ALLOWS SSH ACCESS TO THE EC2 INSTANCES
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_security_group" "example" {
  vpc_id = "${var.vpc_id}"

  # Outbound Everything
  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Inbound SSH from anywhere
  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "${var.name_prefix}-example"
  }
}