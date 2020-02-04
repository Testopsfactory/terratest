terraform {
  source = "..//terraform-with-plan-error"
  extra_arguments "common_vars" {
    arguments = [
      "-var-file=terraform.tfvars"
    ]
  }
}
