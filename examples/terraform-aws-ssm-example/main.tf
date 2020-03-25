provider "aws" {
  region = "us-east-2"
}

resource "aws_iam_role" "example" {
  name = "example"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
        Effect = "Allow"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "example-ssm" {
  role       = aws_iam_role.example.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforSSM"
}

resource "aws_iam_instance_profile" "example" {
  name = "example"
  role = aws_iam_role.example.name
}

resource "aws_instance" "example" {
  ami                         = "ami-04100f1cdba76b497"
  instance_type               = "t2.micro"
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.example.name

  # website::tag::1:: When the instance boots, install the SSM agent.
  user_data = <<-EOF
    #!/bin/bash
    mkdir /tmp/ssm
    curl https://s3.amazonaws.com/ec2-downloads-windows/SSMAgent/latest/debian_amd64/amazon-ssm-agent.deb -o /tmp/ssm/amazon-ssm-agent.deb
    sudo dpkg -i /tmp/ssm/amazon-ssm-agent.deb
    rm -rf /tmp/ssm
    sudo systemctl enable amazon-ssm-agent
  EOF
}

# website::tag::2:: Output the instance's public IP address.
output "instance_id" {
  value = aws_instance.example.id
}
