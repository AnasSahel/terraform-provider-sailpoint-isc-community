# List all available connectors
data "sailpoint_connectors" "all" {}

# List connectors with filtering for Active Directory type connectors
data "sailpoint_connectors" "active_directory" {
  filters = "name sw \"Active Directory\""
}

# List connectors with pagination
data "sailpoint_connectors" "paginated" {
  limit         = 10
  offset        = 0
  include_count = true
}

# List connectors with German locale
data "sailpoint_connectors" "german" {
  locale = "de"
}

# Fetch all connectors using SailPoint pagination (up to 10,000)
data "sailpoint_connectors" "all_paginated" {
  paginate_all = true
}

# Fetch connectors with custom pagination settings
data "sailpoint_connectors" "custom_paginated" {
  paginate_all = true
  max_results  = 5000 # Maximum 5,000 results
  page_size    = 100  # Fetch 100 at a time
  filters      = "status eq \"RELEASED\""
}

# Fetch filtered connectors with pagination
data "sailpoint_connectors" "filtered_paginated" {
  paginate_all = true
  filters      = "directConnect eq true"
  locale       = "en"
}

# Output examples
output "all_connectors" {
  description = "List of all available connectors"
  value       = data.sailpoint_connectors.all.connectors
}

output "connector_count" {
  description = "Total number of connectors"
  value       = length(data.sailpoint_connectors.all.connectors)
}

output "active_directory_connectors" {
  description = "Active Directory type connectors"
  value       = data.sailpoint_connectors.active_directory.connectors
}

output "all_paginated_connectors" {
  description = "All connectors fetched using pagination"
  value       = data.sailpoint_connectors.all_paginated.connectors
}

output "released_connectors_count" {
  description = "Count of released connectors (with custom pagination)"
  value       = length(data.sailpoint_connectors.custom_paginated.connectors)
}

output "direct_connect_connectors" {
  description = "Connectors that support direct connect"
  value = [
    for conn in data.sailpoint_connectors.filtered_paginated.connectors :
    {
      name           = conn.name
      type           = conn.type
      script_name    = conn.script_name
      direct_connect = conn.direct_connect
      features       = conn.features
    }
  ]
}
