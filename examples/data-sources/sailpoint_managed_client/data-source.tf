# Look up a virtual appliance by ID.
data "sailpoint_managed_client" "by_id" {
  id = "REPLACE_WITH_MANAGED_CLIENT_ID"
}

# Or by name (fails if zero or several managed clients match).
data "sailpoint_managed_client" "by_name" {
  name = "Example VA"
}

# Narrow a name lookup to one cluster.
data "sailpoint_managed_client" "in_cluster" {
  name       = "Example VA"
  cluster_id = "REPLACE_WITH_CLUSTER_ID"
}

output "va_status" {
  value = data.sailpoint_managed_client.by_name.status
}
