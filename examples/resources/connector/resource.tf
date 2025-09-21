# Basic custom connector
resource "sailpoint_connector" "basic" {
  name = "My Custom Connector"
  type = "custom-rest"
}

# Advanced custom connector with all options
resource "sailpoint_connector" "advanced" {
  name           = "Advanced Custom Connector"
  type           = "custom-advanced"
  script_name    = "advanced_custom_connector"
  class_name     = "sailpoint.connector.OpenConnectorAdapter"
  direct_connect = true
  file_upload    = false

  # Connector metadata as JSON
  connector_metadata = jsonencode({
    "supportedObjectTypes" = ["account", "group", "entitlement"]
    "displayName"          = "Advanced Custom Connector"
    "description"          = "A custom connector built with Terraform"
    "features"             = ["PROVISIONING", "SYNC_PROVISIONING", "PASSWORD"]
  })

  # Translation properties for localization
  translation_properties = jsonencode({
    "en" = {
      "displayName" = "Advanced Custom Connector"
      "description" = "A custom connector built with Terraform"
    }
    "fr" = {
      "displayName" = "Connecteur Personnalisé Avancé"
      "description" = "Un connecteur personnalisé construit avec Terraform"
    }
  })
}

# Output the connector details
output "basic_connector_id" {
  description = "The ID (script name) of the basic custom connector"
  value       = sailpoint_connector.basic.id
}

output "advanced_connector_script_name" {
  description = "The script name of the advanced custom connector"
  value       = sailpoint_connector.advanced.script_name
}
