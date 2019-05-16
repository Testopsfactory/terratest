 terragrunt = {
     terraform = {
         source = "..//terraform-with-error"
         arguments = [
             "-var-file=terraform.tfvars"
         ]
     }
 }