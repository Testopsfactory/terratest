provider "aws" {
  region = var.region
}

data "aws_iam_policy_document" "example" {
  version = "2012-10-17"

  statement {
    sid = "1"

    actions = [
      "sts:AssumeRole",
    ]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "example" {
  name_prefix        = "example"
  assume_role_policy = data.aws_iam_policy_document.example.json
}

resource "aws_iam_role_policy_attachment" "example-ssm" {
  role       = aws_iam_role.example.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforSSM"
}

resource "aws_iam_instance_profile" "example" {
  name_prefix = "example"
  role        = aws_iam_role.example.name
}

data "aws_ami" "amazon-linux-2" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm*"]
  }
}

resource "aws_instance" "example" {
  ami                         = data.aws_ami.amazon-linux-2.id
  instance_type               = "t2.micro"
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.example.name
}
