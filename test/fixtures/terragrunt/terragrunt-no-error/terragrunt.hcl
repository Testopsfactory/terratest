terraform {
  source = "..//terragrunt-no-error"
  extra_arguments "common_vars" {
    arguments = [
      "-var-file=terraform.tfvars"
    ]
  }
}
