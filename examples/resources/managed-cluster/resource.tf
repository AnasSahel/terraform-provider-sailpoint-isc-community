# Example 1: Basic managed cluster resource
resource "sailpoint_managed_cluster" "basic_example" {
  name        = "Test Cluster"
  description = "A test managed cluster for development"
  type        = "idn"

  configuration = {
    gmt_offset = "-5"
  }
}

# Example 2: Using variables for flexible configuration
variable "cluster_name" {
  description = "Name of the managed cluster"
  type        = string
  default     = "My Cluster"
}

variable "environment" {
  description = "Environment (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "timezone_offset" {
  description = "GMT timezone offset"
  type        = string
  default     = "-5"
}

resource "sailpoint_managed_cluster" "variable_example" {
  name        = var.cluster_name
  description = "Managed cluster for ${var.environment} environment"
  type        = "idn"

  configuration = {
    gmt_offset = var.timezone_offset
    log_level  = var.environment == "prod" ? "WARN" : "DEBUG"
  }
}

# Example 3: Multiple clusters with for_each
variable "clusters" {
  description = "Map of clusters to create"
  type = map(object({
    description = string
    gmt_offset  = string
    log_level   = string
  }))
  default = {
    "dev-cluster" = {
      description = "Development cluster"
      gmt_offset  = "-5"
      log_level   = "DEBUG"
    }
    "staging-cluster" = {
      description = "Staging cluster"
      gmt_offset  = "-5"
      log_level   = "INFO"
    }
  }
}

resource "sailpoint_managed_cluster" "multiple_clusters" {
  for_each = var.clusters

  name        = each.key
  description = each.value.description
  type        = "idn"

  configuration = {
    gmt_offset = each.value.gmt_offset
    log_level  = each.value.log_level
  }
}

# Example 4: Output cluster information
output "cluster_ids" {
  description = "IDs of created managed clusters"
  value = {
    basic         = sailpoint_managed_cluster.basic_example.id
    comprehensive = sailpoint_managed_cluster.comprehensive_example.id
    variable      = sailpoint_managed_cluster.variable_example.id
  }
}

output "cluster_info" {
  description = "Detailed information about the basic cluster"
  value = {
    id                    = sailpoint_managed_cluster.basic_example.id
    name                  = sailpoint_managed_cluster.basic_example.name
    type                  = sailpoint_managed_cluster.basic_example.type
    operational           = sailpoint_managed_cluster.basic_example.operational
    status                = sailpoint_managed_cluster.basic_example.status
    client_type           = sailpoint_managed_cluster.basic_example.client_type
    ccg_version           = sailpoint_managed_cluster.basic_example.ccg_version
    service_count         = sailpoint_managed_cluster.basic_example.service_count
    public_key_thumbprint = sailpoint_managed_cluster.basic_example.public_key_thumbprint
    created_at            = sailpoint_managed_cluster.basic_example.created_at
    updated_at            = sailpoint_managed_cluster.basic_example.updated_at
  }
}

# Example 5: Import existing managed cluster
# To import: terraform import sailpoint_managed_cluster.imported_cluster 2c918085-74f3-4b96-8c31-3c3a7cb8f5e2
resource "sailpoint_managed_cluster" "imported_cluster" {
  name        = "Existing Cluster"
  description = "Previously created cluster now managed by Terraform"
  type        = "idn"

  configuration = {
    gmt_offset = "-5"
  }

  lifecycle {
    # Prevent accidental deletion of important clusters
    prevent_destroy = true
  }
}
