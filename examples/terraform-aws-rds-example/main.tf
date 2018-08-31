resource "aws_vpc" "example" {
  cidr_block = "10.0.0.0/16"
  tags {
    Name = "${var.name}"
  }
}

data "aws_availability_zones" "available" {}


resource "aws_subnet" "example_1" {
  vpc_id     = "${aws_vpc.example.id}"
  cidr_block = "10.0.0.0/24"
  availability_zone = "${data.aws_availability_zones.available.names[0]}"
  tags {
    Name = "${var.name}-1"
  }
}

resource "aws_subnet" "example_2" {
  vpc_id     = "${aws_vpc.example.id}"
  cidr_block = "10.0.1.0/24"
  availability_zone = "${data.aws_availability_zones.available.names[1]}"
  tags {
    Name = "${var.name}-2"
  }
}

resource "aws_db_subnet_group" "example" {
  name       = "${var.name}"
  subnet_ids = ["${aws_subnet.example_1.id}", "${aws_subnet.example_2.id}"]

  tags {
    Name = "${var.name}"
  }
}

resource "aws_db_option_group" "example" {
  name                     = "${var.name}"
  engine_name              = "${var.engine_name}"
  major_engine_version     = "${var.major_engine_version}"

  tags {
    Name = "${var.name}"
  }

  option {
    option_name  = "MARIADB_AUDIT_PLUGIN"

    option_settings {
      name  = "SERVER_AUDIT_EVENTS"
      value = "CONNECT"
    }
  }
}

resource "aws_db_parameter_group" "example" {
  name        = "${var.name}"
  family      = "${var.family}"

  tags {
    Name = "${var.name}"
  }

  parameter {
    name  = "general_log"
    value = "0"
  }

}

resource "aws_db_instance" "example" {
  allocated_storage     = "${var.allocated_storage}"
  engine                = "${var.engine_name}"
  instance_class        = "db.t2.micro"
  license_model         = "${var.license_model}"
  username              = "${var.username}"
  password              = "${var.password}"
  db_subnet_group_name  = "${aws_db_subnet_group.example.id}"
  skip_final_snapshot   = true
  identifier            = "${var.name}"

  tags {
    Name = "${var.name}"
  }

  parameter_group_name = "${aws_db_parameter_group.example.id}"
  option_group_name    = "${aws_db_option_group.example.id}"
}
