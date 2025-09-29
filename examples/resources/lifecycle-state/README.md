# SailPoint Lifecycle State Resource Example

This example demonstrates how to create and manage lifecycle states in SailPoint Identity Security Cloud (ISC) using the `sailpoint_lifecycle_state` resource.

## Usage

**Important**: This resource requires an existing identity profile in your SailPoint ISC tenant. You must provide the identity profile ID through variables.

### Method 1: Using terraform.tfvars

Copy the example file and update with your values:

```bash
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your actual values
```

### Method 2: Using environment variables

```bash
export TF_VAR_identity_profile_id="your-identity-profile-id"
export TF_VAR_access_profile_ids='["access-profile-id-1", "access-profile-id-2"]'
terraform plan
```

### Method 3: Command line

```bash
terraform plan \
  -var="identity_profile_id=your-identity-profile-id" \
  -var='access_profile_ids=["access-profile-id-1"]'
```

## Requirements

- A SailPoint Identity Security Cloud (ISC) tenant
- Valid SailPoint credentials configured in the provider
- An existing identity profile where the lifecycle state will be created
- Optionally, existing access profiles to associate with the lifecycle state

## Finding Required IDs

### Identity Profile ID
1. Navigate to Admin → Identity Profiles in SailPoint ISC
2. Select your identity profile
3. Copy the ID from the URL (e.g., `/admin/identity-profiles/{identity-profile-id}`)

### Access Profile IDs (Optional)
1. Navigate to Admin → Access Profiles
2. Select the access profiles you want to associate
3. Copy their IDs from the URLs or configuration

## Resource Configuration

### Required Attributes
- `identity_profile_id` - The identity profile that will contain this lifecycle state
- `name` - Human-readable display name for the lifecycle state
- `technical_name` - Internal technical identifier (must be unique within the profile)

### Optional Attributes
- `description` - Detailed description of the lifecycle state's purpose
- `enabled` - Whether the state is active (defaults to computed value from API)
- `access_profile_ids` - List of access profiles to automatically grant
- `priority` - Priority level for conflict resolution (lower = higher priority)
- `identity_state` - Broader identity state category (usually "ACTIVE" or null)

### Computed Attributes
- `id` - System-generated unique identifier
- `identity_count` - Current number of identities in this state
- `created` - Creation timestamp
- `modified` - Last modification timestamp

## Lifecycle Management

### Create
```bash
terraform plan
terraform apply
```

### Update
Modify the configuration and run:
```bash
terraform plan
terraform apply
```

**Note**: Some attributes like `name` and `technical_name` require resource replacement when changed.

### Import Existing
Use the provided import script:
```bash
./import.sh <identity_profile_id> <lifecycle_state_id>
```

Or manually:
```bash
terraform import sailpoint_lifecycle_state.example "identity_profile_id:lifecycle_state_id"
```

### Delete
```bash
terraform destroy
```

## Example Use Cases

### 1. Basic Lifecycle State
```hcl
resource "sailpoint_lifecycle_state" "basic" {
  identity_profile_id = var.identity_profile_id
  name                = "Basic Active"
  technical_name      = "basic_active"
  description         = "Standard active state for employees"
  enabled             = true
}
```

### 2. Contractor State with Limited Access
```hcl
resource "sailpoint_lifecycle_state" "contractor" {
  identity_profile_id = var.identity_profile_id
  name                = "Contractor Active"
  technical_name      = "contractor_active"
  description         = "Active state for contractors with limited access"
  enabled             = true
  access_profile_ids  = [var.contractor_access_profile_id]
  priority            = 20
  identity_state      = "ACTIVE"
}
```

### 3. High-Priority Executive State
```hcl
resource "sailpoint_lifecycle_state" "executive" {
  identity_profile_id = var.identity_profile_id
  name                = "Executive Active"
  technical_name      = "executive_active"
  description         = "High-priority state for executive identities"
  enabled             = true
  access_profile_ids  = var.executive_access_profile_ids
  priority            = 1  # Highest priority
  identity_state      = "ACTIVE"
}
```

### 4. Inactive/Disabled State
```hcl
resource "sailpoint_lifecycle_state" "inactive" {
  identity_profile_id = var.identity_profile_id
  name                = "Inactive"
  technical_name      = "inactive"
  description         = "Disabled state with no access"
  enabled             = false
  access_profile_ids  = []  # No access profiles
  priority            = 100
  identity_state      = "INACTIVE"
}
```

## Best Practices

1. **Naming Convention**: Use clear, descriptive names for both `name` and `technical_name`
2. **Priority Management**: Assign priorities strategically (1-10 for high priority, 90-100 for low priority)
3. **Access Profile Management**: Keep access profile associations minimal and focused
4. **Documentation**: Always include meaningful descriptions
5. **Testing**: Test lifecycle states in a sandbox environment before production deployment

## Troubleshooting

### Common Issues

1. **Invalid Identity Profile ID**: Ensure the identity profile exists and the ID is correct
2. **Duplicate Technical Name**: Technical names must be unique within the identity profile
3. **Invalid Access Profile IDs**: Verify all access profile IDs exist and are accessible
4. **Permission Issues**: Ensure your SailPoint credentials have sufficient permissions

### Validation

After creation, verify the lifecycle state:
```bash
# Check the outputs
terraform output lifecycle_state_id
terraform output lifecycle_state_enabled

# Verify in SailPoint ISC UI
# Navigate to Admin → Identity Profiles → [Your Profile] → Lifecycle States
```