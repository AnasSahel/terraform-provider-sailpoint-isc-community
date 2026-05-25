# Correlation config for an Active Directory source.
#
# Correlation determines how incoming account records are matched to known
# identities. The most common pattern is to correlate by a unique account
# attribute (e.g. sAMAccountName) against a unique identity attribute (e.g. uid).
#
# This is a singleton-per-source resource. Create adopts the existing config
# (it always exists for a source) and Delete is a no-op — removing this resource
# from Terraform only removes it from Terraform state; the config persists in ISC.
resource "sailpoint_source_correlation_config" "ad" {
  source_id = sailpoint_source.ad.id

  attribute_assignments = [
    {
      property  = "sAMAccountName"
      value     = "uid"
      operation = "EQ"
    },
  ]
}

# Import an existing correlation config into Terraform state:
#   terraform import sailpoint_source_correlation_config.ad <source-id>
