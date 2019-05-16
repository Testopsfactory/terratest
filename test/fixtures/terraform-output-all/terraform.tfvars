 terragrunt = {
     terraform = {
         source = "..//terraform-output-all"
         arguments = [
             "-var-file=terraform.tfvars"
         ]
     }
 }