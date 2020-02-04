terraform {
  source = "..//foo"
  extra_arguments "common_vars" {
    arguments = [
      "-var-file=terraform.tfvars"
    ]
  }
}
