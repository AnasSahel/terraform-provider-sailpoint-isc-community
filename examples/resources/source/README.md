# SailPoint Source Examples

This directory contains examples for managing SailPoint Identity Security Cloud (ISC) sources using the Terraform provider.

## Resource Examples

### Active Directory Source

The Active Directory example demonstrates how to configure a comprehensive AD source with:

- Domain controller configuration
- SSL/TLS encryption settings
- User and group search filters
- Forest settings for multi-domain environments
- Service account authentication
- Provisioning features

```hcl
resource "sailpoint_source" "active_directory" {
  name        = "Corporate Active Directory"
  description = "Main Active Directory source for employee identities"
  connector   = "active-directory"
  
  owner = jsonencode({
    type = "IDENTITY"
    id   = "2c91808570313110017040b06f344ec9"
    name = "john.doe"
  })
  
  # Optional core attributes
  authoritative     = true
  delete_threshold  = 10
  features          = ["PROVISIONING", "NO_PERMISSIONS_PROVISIONING", "GROUPS_HAVE_MEMBERS"]
  
  # Connector-specific configuration as JSON
  connector_attributes = jsonencode({
    domain_name           = "corp.example.com"
    domain_controller     = "dc1.corp.example.com"
    # ... additional configuration
  })
}
```

### Delimited File Source

The CSV file example shows how to configure a file-based source for bulk imports:

- File format and delimiter configuration
- Column mapping and merging
- Identity attribute specification
- Group handling from CSV columns

```hcl
resource "sailpoint_source" "employee_csv" {
  name        = "Employee CSV Import"
  description = "CSV file source for bulk employee data import"
  connector   = "delimited-file"
  
  owner = jsonencode({
    type = "IDENTITY"
    id   = "2c91808570313110017040b06f344ec9"
    name = "john.doe"
  })
  
  # Connector-specific configuration as JSON
  connector_attributes = jsonencode({
    file               = "employees.csv"
    delimiter          = ","
    has_header         = "true"
    column_names       = "username,firstName,lastName,email,department,title,manager"
    identity_attribute = "username"
  })
}
```

## Data Source Examples

### List All Sources

Retrieve information about all sources in your tenant:

```hcl
data "sailpoint_sources" "all" {}

output "all_source_names" {
  value = [for source in data.sailpoint_sources.all.sources : source.name]
}
```

### Get Specific Source

Retrieve a specific source by ID or name:

```hcl
data "sailpoint_source" "by_id" {
  id = "2c91808570313110017040b06f344ec9"
}

data "sailpoint_source" "by_name" {
  name = "Corporate Active Directory"
}
```

## Import Existing Sources

Use the included import script to bring existing SailPoint sources under Terraform management:

```bash
# Make the script executable (if not already)
chmod +x import.sh

# Import a source
./import.sh 2c91808570313110017040b06f344ec9 sailpoint_source.imported_ad
```

The script will:
1. Import the source into your Terraform state
2. Provide guidance on next steps
3. Show example configuration structure

### Manual Import

You can also import manually using terraform import:

```bash
terraform import sailpoint_source.my_source 2c91808570313110017040b06f344ec9
```

## Configuration Reference

### Required Fields

- `name` - Unique name for the source
- `description` - Description of the source
- `connector` - Connector type (e.g., "active-directory", "delimited-file")
- `owner` - JSON-encoded owner object with type, id, and name

### Optional Fields

#### Core Attributes
- `type` - Type of system being managed (computed)
- `connector_class` - Java class implementing the connector (computed)
- `connection_type` - Type of connection (direct or file)
- `authoritative` - Whether source is authoritative for identities
- `cluster` - JSON-encoded cluster assignment (for VA sources)

#### Configuration
- `connector_attributes` - JSON-encoded connector-specific configuration (sensitive)
- `delete_threshold` - Account deletion threshold (0-100)
- `features` - List of enabled features (e.g., "PROVISIONING")

#### Management & Correlation
- `management_workgroup` - JSON-encoded workgroup for source management
- `account_correlation_config` - JSON-encoded correlation configuration
- `account_correlation_rule` - JSON-encoded correlation rule
- `manager_correlation_rule` - JSON-encoded manager correlation rule
- `manager_correlation_mapping` - JSON-encoded correlation mapping

#### Provisioning
- `before_provisioning_rule` - JSON-encoded pre-provisioning rule
- `password_policies` - JSON-encoded password policy list

### Common Configuration Options

#### Active Directory (`connector_attributes` JSON)
- `domain_name` - AD domain name
- `domain_controller` - Primary domain controller
- `forest_settings` - Multi-domain forest configuration array
- `user_search_filter` - LDAP filter for users
- `group_search_filter` - LDAP filter for groups
- `authorization_type` - Authentication method (simple, etc.)
- `use_tls` - Enable TLS encryption
- `account_username` - Service account username
- `account_password` - Service account password (use variables!)

#### Delimited File (`connector_attributes` JSON)
- `file` - Path to the CSV file
- `delimiter` - Field delimiter character
- `has_header` - Whether file has header row
- `column_names` - Comma-separated column names
- `identity_attribute` - Primary identity column
- `group_column_name` - Column containing group memberships
- `merge_columns` - Object for merging multiple columns

## Best Practices

1. **Security**: Store sensitive values (passwords, keys) in Terraform variables marked as sensitive
2. **Validation**: Use `terraform plan` to review changes before applying
3. **State Management**: Keep Terraform state secure and backed up
4. **Naming**: Use consistent naming conventions for sources
5. **Documentation**: Comment your configuration for team collaboration

## Troubleshooting

### Import Issues
- Verify source ID exists in SailPoint
- Check SailPoint API credentials
- Ensure proper network connectivity

### Configuration Errors
- Validate JSON-encoded fields are properly formatted
- Check required connector-specific configuration options
- Review SailPoint documentation for connector requirements

### Plan/Apply Issues
- Run `terraform plan` to identify configuration drift
- Check for conflicting manual changes in SailPoint UI
- Verify all required fields are present

## Additional Resources

- [SailPoint ISC API Documentation](https://developer.sailpoint.com/idn/api)
- [Terraform Provider Documentation](../../docs/)
- [SailPoint Connector Guide](https://documentation.sailpoint.com/connectors/)