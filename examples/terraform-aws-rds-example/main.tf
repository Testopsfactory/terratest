data "aws_vpc" "default" {
  default = true
}

data "aws_availability_zones" "available" {}

data "aws_subnet" "example_1" {
  vpc_id     = "${data.aws_vpc.default.id}"
  availability_zone = "${data.aws_availability_zones.available.names[0]}"
}

data "aws_subnet" "example_2" {
  vpc_id     = "${data.aws_vpc.default.id}"
  availability_zone = "${data.aws_availability_zones.available.names[1]}"
}

resource "aws_db_subnet_group" "example" {
  name       = "${var.name}"
  subnet_ids = ["${data.aws_subnet.example_1.id}", "${data.aws_subnet.example_2.id}"]

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
  engine_version        = "${var.engine_version}"
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
