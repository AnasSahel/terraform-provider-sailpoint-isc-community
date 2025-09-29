# SailPoint Lifecycle State List Data Source Example

This example demonstrates how to use the `sailpoint_lifecycle_state_list` data source to retrieve information about all lifecycle states within an identity profile in SailPoint Identity Security Cloud (ISC).

## Usage

**Important**: This data source requires an existing identity profile in your SailPoint ISC tenant. You must provide the identity profile ID through variables.

### Method 1: Using terraform.tfvars

Create a `terraform.tfvars` file:

```hcl
identity_profile_id = "your-identity-profile-id"
```

### Method 2: Using environment variables

```bash
export TF_VAR_identity_profile_id="your-identity-profile-id"
terraform plan
```

### Method 3: Command line

```bash
terraform plan -var="identity_profile_id=your-identity-profile-id"
```

## Requirements

- A SailPoint Identity Security Cloud (ISC) tenant
- Valid SailPoint credentials configured in the provider
- An existing identity profile with lifecycle states

## Finding the Identity Profile ID

To find the required identity profile ID in SailPoint ISC:

1. Navigate to Admin → Identity Profiles
2. Select your identity profile
3. Copy the ID from the URL (e.g., `/admin/identity-profiles/{identity-profile-id}`)

## Available Outputs

The data source provides comprehensive information about all lifecycle states:

### Individual Lifecycle State Fields
- `id` - System-generated unique identifier for each lifecycle state
- `name` - Human-readable display name (e.g., 'Active', 'Inactive', 'Terminated')
- `technical_name` - Internal technical name used by SailPoint systems
- `enabled` - Boolean indicating if the state is currently active
- `description` - Detailed description of the lifecycle state's purpose
- `identity_count` - Current number of identities in this lifecycle state
- `access_profile_ids` - List of access profiles automatically granted in this state
- `priority` - Priority level for conflict resolution between states
- `created` - ISO 8601 timestamp when the state was created
- `modified` - ISO 8601 timestamp when the state was last modified

### Example Outputs Provided
- `lifecycle_states_count` - Total number of lifecycle states
- `lifecycle_state_names` - Array of all lifecycle state names
- `enabled_lifecycle_states` - Filtered list of only enabled states
- `lifecycle_states_access_summary` - Summary with access profile counts and priorities

## Example Use Cases

1. **Audit & Compliance**: Review all lifecycle states and their configurations
2. **Monitoring**: Track identity distribution across different lifecycle states
3. **Access Management**: Analyze which access profiles are granted in each state
4. **Reporting**: Generate comprehensive reports on identity profile configurations
5. **Automation**: Use lifecycle state data to trigger other Terraform resources
6. **Capacity Planning**: Understand identity counts and access patterns

## Advanced Usage Examples

### Filter by Identity Count
```hcl
locals {
  high_usage_states = [
    for state in data.sailpoint_lifecycle_state_list.example.lifecycle_state_list : state
    if state.identity_count > 100
  ]
}
```

### Group by Priority
```hcl
locals {
  states_by_priority = {
    for state in data.sailpoint_lifecycle_state_list.example.lifecycle_state_list : 
    state.priority => state.name...
  }
}
```

## Testing

When running acceptance tests, use environment variables:

```bash
export TF_VAR_identity_profile_id="your-test-identity-profile-id"
make testacc
```