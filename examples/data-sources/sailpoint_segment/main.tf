# Look up an existing segment by ID
data "sailpoint_segment" "example" {
  id = "2c91808a7813090a017ecccc00000001"
}

output "segment_name" {
  value = data.sailpoint_segment.example.name
}

output "segment_active" {
  value = data.sailpoint_segment.example.active
}

output "segment_visibility_criteria" {
  value = data.sailpoint_segment.example.visibility_criteria
}
