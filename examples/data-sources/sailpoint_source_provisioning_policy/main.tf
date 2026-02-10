# Look up the CREATE provisioning policy for a source
data "sailpoint_source_provisioning_policy" "create" {
  source_id  = "2c91808a7813090a017813467c744b01"
  usage_type = "CREATE"
}

# Output the policy details
output "policy_name" {
  value = data.sailpoint_source_provisioning_policy.create.name
}

output "policy_fields" {
  value = data.sailpoint_source_provisioning_policy.create.fields
}

# Look up the UPDATE provisioning policy for the same source
data "sailpoint_source_provisioning_policy" "update" {
  source_id  = "2c91808a7813090a017813467c744b01"
  usage_type = "UPDATE"
}
