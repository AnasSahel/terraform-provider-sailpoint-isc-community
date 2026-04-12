# Look up an existing access profile by ID
data "sailpoint_access_profile" "example" {
  id = "REPLACE_WITH_ACCESS_PROFILE_ID"
}

output "access_profile_name" {
  value = data.sailpoint_access_profile.example.name
}

output "access_profile_entitlements" {
  value = data.sailpoint_access_profile.example.entitlements
}
