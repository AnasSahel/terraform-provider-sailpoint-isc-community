# SailPoint Lifecycle State Data Source Example

This example demonstrates how to use the `sailpoint_lifecycle_state` data source to retrieve information about a specific lifecycle state within an identity profile.

## Usage

**Important**: This data source requires existing lifecycle states and identity profiles in your SailPoint ISC tenant. You must provide the actual IDs through variables.

### Method 1: Using terraform.tfvars

Create a `terraform.tfvars` file:

```hcl
lifecycle_state_id  = "your-lifecycle-state-id"
identity_profile_id = "your-identity-profile-id"
```

### Method 2: Using environment variables

```bash
export TF_VAR_lifecycle_state_id="your-lifecycle-state-id"
export TF_VAR_identity_profile_id="your-identity-profile-id"
terraform plan
```

### Method 3: Command line

```bash
terraform plan -var="lifecycle_state_id=your-lifecycle-state-id" -var="identity_profile_id=your-identity-profile-id"
```

## Requirements

- A SailPoint Identity Security Cloud (ISC) tenant
- Valid SailPoint credentials configured in the provider
- An existing identity profile and lifecycle state
- The lifecycle state must exist within the specified identity profile

## Finding Required IDs

To find the required IDs in SailPoint ISC:

1. **Identity Profile ID**: 
   - Navigate to Admin → Identity Profiles
   - Select your identity profile
   - Copy the ID from the URL (e.g., `/admin/identity-profiles/{identity-profile-id}`)

2. **Lifecycle State ID**: 
   - Within the identity profile, go to the Lifecycle States tab
   - Select the desired lifecycle state
   - Copy the ID from the URL or use the SailPoint API to list lifecycle states

## Available Outputs

The data source provides comprehensive information about the lifecycle state:

- `name` - Display name of the lifecycle state
- `technical_name` - Internal technical name used by SailPoint
- `enabled` - Whether the state is currently enabled and active
- `description` - Detailed description of the lifecycle state's purpose
- `identity_count` - Current number of identities in this lifecycle state
- `access_profile_ids` - List of access profiles automatically granted in this state
- `priority` - Priority level used for conflict resolution between states
- `created` - ISO 8601 timestamp when the state was created
- `modified` - ISO 8601 timestamp when the state was last modified

## Example Use Cases

1. **Monitoring**: Track how many identities are in specific lifecycle states
2. **Access Review**: Verify which access profiles are automatically granted
3. **Compliance**: Ensure lifecycle states are properly configured and enabled
4. **Automation**: Use lifecycle state data to trigger other Terraform resources
5. **Reporting**: Generate reports on identity distribution across lifecycle states

## Testing

When running acceptance tests, use environment variables to avoid hardcoding real IDs:

```bash
export TF_VAR_lifecycle_state_id="your-test-lifecycle-state-id"
export TF_VAR_identity_profile_id="your-test-identity-profile-id"
make testacc
```