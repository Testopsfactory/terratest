resource "random_string" "default" {
  length = 3  
  lower = true
  number = false
  special = false
}

resource "azurerm_resource_group" "main" {
  name     =  format("%s-%s-%s", "terratest", lower(random_string.default.result), "keyvault")
  location = var.location
}
