# Look up an application by ID.
data "sailpoint_application" "by_id" {
  id = "REPLACE_WITH_APPLICATION_ID"
}

# Or by name (fails if zero or several applications match).
data "sailpoint_application" "by_name" {
  name = "Example Application"
}
