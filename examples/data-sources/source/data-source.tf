terraform {
  required_providers {
    sailpoint = {
      source = "anasSahel/sailpoint-isc-community"
    }
  }
}

provider "sailpoint" {
  # Configuration options
}

# Get a specific source by ID
data "sailpoint_source" "by_id" {
  id = "2c91808570313110017040b06f344ec9"
}

# Get a specific source by name
data "sailpoint_source" "by_name" {
  name = "Corporate Active Directory"
}

# Use source data in other resources
resource "sailpoint_transform" "source_based_transform" {
  name = "Get Source Name Transform"
  type = "static"

  attributes = jsonencode({
    value = data.sailpoint_source.by_name.name
  })
}

# Output source details
output "source_details" {
  description = "Details of the retrieved source"
  value = {
    id          = data.sailpoint_source.by_id.id
    name        = data.sailpoint_source.by_id.name
    connector   = data.sailpoint_source.by_id.connector
    description = data.sailpoint_source.by_id.description
    owner_name  = jsondecode(data.sailpoint_source.by_id.owner).name
  }
}
