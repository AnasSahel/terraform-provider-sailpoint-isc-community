# Test configuration for managed cluster resource and data source
# This file is used for testing and validation purposes

terraform {
  required_providers {
    sailpoint = {
      source = "registry.terraform.io/providers/AnasSahel/sailpoint-isc-community"
    }
  }
}

# Configure the SailPoint provider
provider "sailpoint" {
  # Configuration will be picked up from environment variables:
  # SAILPOINT_BASE_URL, SAILPOINT_CLIENT_ID, SAILPOINT_CLIENT_SECRET
}

# Test resource creation
resource "sailpoint_managed_cluster" "test" {
  name        = "Terraform Test Cluster"
  description = "Test cluster created by Terraform"
  type        = "idn"

  configuration = {
    gmt_offset  = "-05:00"
    environment = "test"
    log_level   = "DEBUG"
  }
}

# Test data source lookup by ID (using the created resource)
data "sailpoint_managed_cluster" "test_by_id" {
  id = sailpoint_managed_cluster.test.id
}

# Test data source lookup by name  
data "sailpoint_managed_cluster" "test_by_name" {
  name = sailpoint_managed_cluster.test.name
}

# Outputs for verification
output "test_cluster_info" {
  description = "Information about the test cluster"
  value = {
    id                  = sailpoint_managed_cluster.test.id
    name                = sailpoint_managed_cluster.test.name
    operational         = sailpoint_managed_cluster.test.operational
    status              = sailpoint_managed_cluster.test.status
    data_source_by_id   = data.sailpoint_managed_cluster.test_by_id.id
    data_source_by_name = data.sailpoint_managed_cluster.test_by_name.id
  }
}
