 terragrunt = {
     terraform = {
         source = "..//terraform-out-map"
         arguments = [
             "-var-file=terraform.tfvars"
         ]
     }
 }