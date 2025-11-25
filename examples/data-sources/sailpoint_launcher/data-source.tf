# Example: Read an existing launcher by ID
data "sailpoint_launcher" "existing_launcher" {
  id = "2c91808a7b5c3e1d017b5c4a8f6d0002"
}

# Use the launcher data
output "launcher_name" {
  value = data.sailpoint_launcher.existing_launcher.name
}

output "launcher_description" {
  value = data.sailpoint_launcher.existing_launcher.description
}

output "launcher_disabled" {
  value = data.sailpoint_launcher.existing_launcher.disabled
}

output "launcher_reference" {
  value = data.sailpoint_launcher.existing_launcher.reference
}

output "launcher_owner" {
  value = data.sailpoint_launcher.existing_launcher.owner
}

# Example: Use launcher data to create a similar launcher
resource "sailpoint_launcher" "cloned_launcher" {
  name        = "${data.sailpoint_launcher.existing_launcher.name} - Copy"
  description = "Cloned from ${data.sailpoint_launcher.existing_launcher.name}"
  type        = data.sailpoint_launcher.existing_launcher.type
  disabled    = true  # Start disabled for safety

  reference = data.sailpoint_launcher.existing_launcher.reference
  config    = data.sailpoint_launcher.existing_launcher.config
}
