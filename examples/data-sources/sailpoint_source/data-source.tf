# Look up a source by its name, which stays the same across tenants.
data "sailpoint_source" "active_directory" {
  name = "Active Directory"
}

# Or look it up by ID.
data "sailpoint_source" "by_id" {
  id = "2c9180835d2e5168015d32f890ca1581"
}

output "active_directory_source_id" {
  value = data.sailpoint_source.active_directory.id
}
