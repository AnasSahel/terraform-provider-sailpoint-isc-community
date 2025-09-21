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
    id           = data.sailpoint_source.by_id.id
    name         = data.sailpoint_source.by_id.name
    connector    = data.sailpoint_source.by_id.connector
    description  = data.sailpoint_source.by_id.description
    owner_name   = data.sailpoint_source.by_id.owner.name
    owner_id     = data.sailpoint_source.by_id.owner.id
    owner_type   = data.sailpoint_source.by_id.owner.type
    cluster_name = data.sailpoint_source.by_id.cluster != null ? data.sailpoint_source.by_id.cluster.name : null
    cluster_id   = data.sailpoint_source.by_id.cluster != null ? data.sailpoint_source.by_id.cluster.id : null
    cluster_type = data.sailpoint_source.by_id.cluster != null ? data.sailpoint_source.by_id.cluster.type : null
  }
}
