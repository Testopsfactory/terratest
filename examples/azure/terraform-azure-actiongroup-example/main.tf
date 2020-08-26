# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE APP SERVICE PLAN
# This is an example of how to deploy an Azure App Service Plan
# ---------------------------------------------------------------------------------------------------------------------

# ------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
# ------------------------------------------------------------------------------

provider "azurerm" {
  features {}
  skip_provider_registration = true
  # To understand why ^ is here, see https://github.com/hashicorp/terraform/issues/18180
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "rg" {
  name     = var.resourceGroupName
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE APP SERVICE PLAN
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_monitor_action_group" "actionGroup" {
  name                = var.appName
  resource_group_name = azurerm_resource_group.rg.name
  short_name          = var.shortName
  tags                = azurerm_resource_group.rg.tags

  dynamic "email_receiver" {
    for_each = var.enableEmail ? ["email_receiver"] : []
    content {
      name                    = var.emailName
      email_address           = var.emailAddress
      use_common_alert_schema = true
    }
  }

  dynamic "sms_receiver" {
    for_each = var.enableSMS ? ["sms_receiver"] : []
    content {
      name         = var.smsName
      country_code = var.smsCountryCode
      phone_number = var.smsPhoneNumber
    }
  }

  dynamic "webhook_receiver" {
    for_each = var.enableWebHook ? ["webhook_receiver"] : []
    content {
      name                    = var.webhookName
      service_uri             = var.webhookServiceUri
      use_common_alert_schema = true
    }
  }

}