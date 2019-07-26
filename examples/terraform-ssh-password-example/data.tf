# ---------------------------------------------------------------------------------------------------------------------
# SET UP A TEMPLATE AROUND THE USER DATA SCRIPT
# ---------------------------------------------------------------------------------------------------------------------

data "template_file" "user_data" {
  template = "${file("${path.module}/user_data.sh")}"

  vars = {
    terratest_password = "${var.terratest_password}"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# LOOK UP THE LATEST UBUNTU AMI
# ---------------------------------------------------------------------------------------------------------------------

data "aws_ami" "ubuntu" {
  most_recent = "true"
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  filter {
    name   = "image-type"
    values = ["machine"]
  }

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-*-amd64-server-*"]
  }
}
