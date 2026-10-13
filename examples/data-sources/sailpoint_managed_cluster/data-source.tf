# Look up a managed cluster (virtual appliance cluster) by name.
data "sailpoint_managed_cluster" "by_name" {
  name = "Example Cluster"
}

# Or by ID.
data "sailpoint_managed_cluster" "by_id" {
  id = "REPLACE_WITH_CLUSTER_ID"
}

# Typical use: reference the cluster from a source instead of hardcoding its ID.
# resource "sailpoint_source" "example" {
#   ...
#   cluster = {
#     type = "CLUSTER"
#     id   = data.sailpoint_managed_cluster.by_name.id
#   }
# }
