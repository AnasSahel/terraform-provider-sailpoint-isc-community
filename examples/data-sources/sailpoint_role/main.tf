# Look up an existing role by ID
data "sailpoint_role" "example" {
  id = "REPLACE_WITH_ROLE_ID"
}

output "role_name" {
  value = data.sailpoint_role.example.name
}

output "role_access_profiles" {
  value = data.sailpoint_role.example.access_profiles
}

output "role_membership" {
  value = data.sailpoint_role.example.membership
}
