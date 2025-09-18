# SailPoint Managed Cluster Examples

This directory contains examples for using the SailPoint managed cluster resource and data source.

## Files

- **resource.tf** - Comprehensive examples of the `sailpoint_managed_cluster` resource
- **import.sh** - Script for importing existing managed clusters into Terraform
- **../data-sources/managed-cluster/data-source.tf** - Examples of the `sailpoint_managed_cluster` data source

## Quick Start

### Creating a Basic Managed Cluster

```terraform
resource "sailpoint_managed_cluster" "example" {
  name        = "My Test Cluster"
  description = "A test managed cluster for development"
  type        = "idn"
  
  configuration = {
    gmt_offset = "-5"
    log_level  = "INFO"
  }
}
```

### Looking Up an Existing Cluster

```terraform
# By ID
data "sailpoint_managed_cluster" "by_id" {
  id = "2c918085-74f3-4b96-8c31-3c3a7cb8f5e2"
}

# By name (uses optimized server-side filtering when available)
data "sailpoint_managed_cluster" "by_name" {
  name = "Production Cluster"
}
```

## Configuration Options

The `configuration` map supports various key-value pairs for cluster settings. All keys should use snake_case format and will be automatically converted to camelCase for the SailPoint API.

Common configuration keys include:

- **gmt_offset** - Timezone offset as integer (e.g., "-5" for UTC-5)
- **region** - AWS/Azure region (e.g., "us-east-1")
- **log_level** - Logging level (DEBUG, INFO, WARN, ERROR)
- **log_retention_days** - Log retention period
- **max_connections** - Maximum number of connections
- **connection_timeout** - Connection timeout in milliseconds
- **ssl_enabled** - Enable SSL/TLS (true/false)
- **environment** - Environment tag (dev, staging, prod)

## Import Instructions

To import an existing managed cluster:

1. Find your cluster ID from the SailPoint ISC admin interface
2. Edit the `import.sh` script with your cluster ID
3. Run the script: `./import.sh`
4. Create a corresponding resource block in your Terraform configuration
5. Run `terraform plan` to see any configuration drift

## Authentication

Ensure you have the following environment variables set:

```bash
export SAILPOINT_BASE_URL="https://[tenant].api.identitynow.com"
export SAILPOINT_CLIENT_ID="your-client-id"
export SAILPOINT_CLIENT_SECRET="your-client-secret"
```

Or configure them in your provider block:

```terraform
provider "sailpoint" {
  base_url      = "https://[tenant].api.identitynow.com"
  client_id     = "your-client-id"  
  client_secret = "your-client-secret"
}
```

## Performance Optimizations

The managed cluster data source uses intelligent filtering to optimize performance:

- **Server-side filtering**: When looking up clusters by name, the data source first attempts to use SailPoint's server-side filtering (`name eq "cluster-name"`) for faster results
- **Automatic fallback**: If server-side filtering is not supported for the name field, it automatically falls back to client-side filtering by retrieving all clusters and filtering locally
- **Logging**: The provider logs when fallback occurs, helping with debugging and API compatibility

## Best Practices

1. **Use descriptive names** - Cluster names should clearly indicate their purpose
2. **Tag with environment** - Use configuration keys to tag clusters by environment
3. **Prevent accidental deletion** - Use `prevent_destroy = true` for production clusters
4. **Use variables** - Parameterize common settings for reusability
5. **Monitor status** - Use the data source to check cluster operational status
6. **Prefer ID lookups** - When possible, use cluster ID for lookups as they're always more efficient

## Troubleshooting

- **Authentication errors**: Verify your SailPoint credentials are correct
- **Cluster not found**: Check the cluster ID/name exists in SailPoint ISC
- **Configuration drift**: Some fields are managed by SailPoint and may show as changes
- **Import failures**: Ensure the cluster isn't already managed by another Terraform state

For more examples, see the files in this directory.