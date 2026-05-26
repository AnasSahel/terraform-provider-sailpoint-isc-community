# Look up a governance group by ID
data "sailpoint_governance_group" "by_id" {
  id = "REPLACE_WITH_GOVERNANCE_GROUP_ID"
}

# Look up a governance group by name using an ISC filter expression
data "sailpoint_governance_group" "by_name" {
  filters = "name eq \"Access Approvers\""
}

# Use the resolved ID in an approval scheme
# (example usage — the approval_schemes attribute will be added to sailpoint_role
# in a future enhancement)
output "approvers_group_id" {
  value = data.sailpoint_governance_group.by_name.id
}
