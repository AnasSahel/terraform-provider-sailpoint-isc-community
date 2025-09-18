# Example 1: Look up managed cluster by ID
data "sailpoint_managed_cluster" "by_id" {
  id = "2c918085-74f3-4b96-8c31-3c3a7cb8f5e2"
}

# Example 2: Look up managed cluster by name
# Uses server-side filtering for better performance when available,
# automatically falls back to client-side filtering if not supported
data "sailpoint_managed_cluster" "by_name" {
  name = "My Test Cluster"
}

# Example 3: Using the data source output in other resources
resource "sailpoint_some_other_resource" "example" {
  managed_cluster_id = data.sailpoint_managed_cluster.by_name.id

  # Access other attributes from the data source
  cluster_name   = data.sailpoint_managed_cluster.by_name.name
  cluster_type   = data.sailpoint_managed_cluster.by_name.type
  is_operational = data.sailpoint_managed_cluster.by_name.operational
}

# Output examples to show available attributes
output "cluster_info" {
  description = "Information about the managed cluster"
  value = {
    id                    = data.sailpoint_managed_cluster.by_name.id
    name                  = data.sailpoint_managed_cluster.by_name.name
    description           = data.sailpoint_managed_cluster.by_name.description
    type                  = data.sailpoint_managed_cluster.by_name.type
    pod                   = data.sailpoint_managed_cluster.by_name.pod
    org                   = data.sailpoint_managed_cluster.by_name.org
    client_type           = data.sailpoint_managed_cluster.by_name.client_type
    ccg_version           = data.sailpoint_managed_cluster.by_name.ccg_version
    operational           = data.sailpoint_managed_cluster.by_name.operational
    status                = data.sailpoint_managed_cluster.by_name.status
    service_count         = data.sailpoint_managed_cluster.by_name.service_count
    client_ids            = data.sailpoint_managed_cluster.by_name.client_ids
    configuration         = data.sailpoint_managed_cluster.by_name.configuration
    public_key_thumbprint = data.sailpoint_managed_cluster.by_name.public_key_thumbprint
    created_at            = data.sailpoint_managed_cluster.by_name.created_at
    updated_at            = data.sailpoint_managed_cluster.by_name.updated_at
  }
}

# Example 4: Conditional lookup (either by ID or name)
data "sailpoint_managed_cluster" "conditional" {
  # Use either id OR name, not both
  id   = var.cluster_id != "" ? var.cluster_id : null
  name = var.cluster_id == "" ? var.cluster_name : null
}
