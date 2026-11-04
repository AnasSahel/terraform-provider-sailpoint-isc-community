data "sailpoint_role_dimension" "london" {
  role_id = "REPLACE_WITH_DIMENSIONAL_ROLE_ID"
  id      = "REPLACE_WITH_DIMENSION_ID"
}

output "london_access_profiles" {
  value = data.sailpoint_role_dimension.london.access_profiles
}

output "london_membership" {
  value = data.sailpoint_role_dimension.london.membership
}
