terraform {
  source = "..//terragrunt-with-error"
  extra_arguments "common_vars" {
    arguments = [
      "-var-file=terraform.tfvars"
    ]
  }
}
