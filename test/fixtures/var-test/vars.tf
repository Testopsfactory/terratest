variable "string" {
  description = "An input to check that we handle string variables correctly"
}

variable "boolean" {
  description = "An input to check that we handle boolean variables correctly"
}

variable "int" {
  description = "An input to check that we handle int variables correctly"
}

variable "map" {
  description = "An input to check that we handle map variables correctly"
  type = "map"
  # If we don't specify a default, we get the error "variable map should be type map, got list". For more info, see:
  # https://github.com/hashicorp/terraform/issues/8057
  default = {}
}

variable "list" {
  description = "An input to check that we handle list variables correctly"
  type = "list"
}


