# Identity Attribute Resource Examples# Identity Attribute Resource Examples# Identity Attribute Resource



This directory contains comprehensive examples for using the `sailpoint_identity_attribute` resource to manage SailPoint Identity Security Cloud identity attributes.



## OverviewThis directory contains comprehensive examples for using the `sailpoint_identity_attribute` resource to manage SailPoint Identity Security Cloud identity attributes.This resource allows you to create and manage identity attributes in SailPoint Identity Security Cloud (ISC).



Identity attributes in SailPoint ISC define the properties and characteristics of identities in your system. This Terraform resource allows you to create, update, and manage custom identity attributes programmatically.



## Files in this Directory## Overview## Example Usage



- **`resource.tf`** - Basic examples covering fundamental usage patterns

- **`advanced-examples.tf`** - Real-world business scenarios and specialized use cases

- **`import-examples.tf`** - Comprehensive import scenarios and troubleshootingIdentity attributes in SailPoint ISC define the properties and characteristics of identities in your system. This Terraform resource allows you to create, update, and manage custom identity attributes programmatically.### Minimal Identity Attribute (only name required)

- **`import-workflow.tf`** - Complete import-first workflow example

- **`integration-patterns.tf`** - Advanced integration patterns with other systems

- **`import.sh`** - Interactive import script with error handling

- **`import_template.sh`** - Template for bulk import operations## Files in this Directory```hcl



## Quick Startresource "sailpoint_identity_attribute" "minimal" {



### Creating New Attributes- **`resource.tf`** - Basic examples covering fundamental usage patterns  name = "costCenter"

```hcl

resource "sailpoint_identity_attribute" "simple" {- **`advanced-examples.tf`** - Real-world business scenarios and specialized use cases  # All other fields have sensible defaults:

  name = "myAttribute"

  # All other properties have sensible defaults- **`import-examples.tf`** - Examples for importing existing identity attributes  # display_name = "costCenter" (defaults to name)

}

```- **`integration-patterns.tf`** - Advanced integration patterns with other systems  # type = "string" (default)



### Importing Existing Attributes  # multi = false (default)

```bash

# Import an existing attribute## Quick Start  # searchable = false (default)

terraform import sailpoint_identity_attribute.cost_center "costCenter"

}

# Verify the import

terraform plan### Minimal Configuration```

```

```hcl

## Import Functionality

resource "sailpoint_identity_attribute" "simple" {### Basic Identity Attribute with explicit values

The SailPoint ISC provider supports importing existing identity attributes into Terraform management. This is essential for managing existing SailPoint deployments.

  name = "myAttribute"

### Import Command

```bash  # All other properties have sensible defaults```hcl

terraform import sailpoint_identity_attribute.<resource_name> "<attribute_name>"

```}resource "sailpoint_identity_attribute" "cost_center" {



### Import Examples```  name         = "costCenter"

```bash

# Import common business attributes  display_name = "Cost Center"

terraform import sailpoint_identity_attribute.employee_id "employeeId"

terraform import sailpoint_identity_attribute.cost_center "costCenter"### Complete Configuration  type         = "string"

terraform import sailpoint_identity_attribute.department "department"

terraform import sailpoint_identity_attribute.manager "manager"```hcl  multi        = false

```

resource "sailpoint_identity_attribute" "complete" {  searchable   = true

### Import Workflow

  name         = "employeeId"}

1. **Discovery**: Use data sources to find existing attributes

   ```hcl  display_name = "Employee ID"```

   data "sailpoint_identity_attribute_list" "existing" {

     include_system = false  type         = "string"

   }

   ```  multi        = false### Identity Attribute with Rule Source



2. **Configuration**: Create resource configurations matching existing state  searchable   = true

   ```hcl

   resource "sailpoint_identity_attribute" "employee_id" {}```hcl

     name         = "employeeId"

     display_name = "Employee ID"```resource "sailpoint_identity_attribute" "department_code" {

     type         = "string"

     searchable   = true  name         = "departmentCode"

     

     lifecycle {## Attribute Properties  display_name = "Department Code"

       prevent_destroy = true  # Safety during import phase

     }  type         = "string"

   }

   ```### Required Properties  multi        = false



3. **Import**: Bring the resource under Terraform management- **`name`** - The unique name of the identity attribute (cannot be changed after creation)  searchable   = true

   ```bash

   terraform import sailpoint_identity_attribute.employee_id "employeeId"

   ```

### Optional Properties  sources {

4. **Validation**: Verify no configuration drift

   ```bash- **`display_name`** - Human-readable display name (defaults to `name`)    type = "rule"

   terraform plan  # Should show no changes

   ```- **`type`** - Data type: `"string"`, `"boolean"`, `"int"`, `"date"` (defaults to `"string"`)    properties = jsonencode({



5. **Management**: Gradually enable full Terraform management- **`multi`** - Whether the attribute can have multiple values (defaults to `false`)      ruleType = "IdentityAttribute"

   ```hcl

   # Remove prevent_destroy when ready for full management- **`searchable`** - Whether the attribute is searchable in the UI (defaults to `false`)      ruleName = "Department Code Mapping Rule"

   lifecycle {

     # prevent_destroy = true  # Remove this line- **`standard`** - Whether this is a standard SailPoint attribute (computed, read-only)    })

   }

   ```- **`system`** - Whether this is a system attribute (computed, read-only)  }



## Attribute Properties}



### Required Properties### Sources Configuration```

- **`name`** - The unique name of the identity attribute (cannot be changed after creation)

The `sources` attribute defines how the identity attribute gets its values. This is an advanced feature that would typically be configured after the attribute is created:

### Optional Properties

- **`display_name`** - Human-readable display name (defaults to `name`)### Multi-Value Identity Attribute

- **`type`** - Data type: `"string"`, `"boolean"`, `"int"`, `"date"` (defaults to `"string"`)

- **`multi`** - Whether the attribute can have multiple values (defaults to `false`)```hcl

- **`searchable`** - Whether the attribute is searchable in the UI (defaults to `false`)

- **`standard`** - Whether this is a standard SailPoint attribute (computed, read-only)# Note: Sources configuration may require additional setup in SailPoint```hcl

- **`system`** - Whether this is a system attribute (computed, read-only)

# Refer to SailPoint documentation for specific source configuration requirementsresource "sailpoint_identity_attribute" "skill_tags" {

## Data Types

```  name         = "skillTags"

| Type | Description | Example Values |

|------|-------------|----------------|  display_name = "Skill Tags"

| `string` | Text values | `"John Doe"`, `"IT Department"` |

| `boolean` | True/false values | `true`, `false` |## Data Types  type         = "string"

| `int` | Integer numbers | `42`, `1000` |

| `date` | Date values | `"2023-12-25"` |  multi        = true



## Import Best Practices| Type | Description | Example Values |  searchable   = true



### 1. Discovery First|------|-------------|----------------|

Always discover existing attributes before importing:

```hcl| `string` | Text values | `"John Doe"`, `"IT Department"` |  sources {

data "sailpoint_identity_attribute_list" "discovery" {

  include_system = true| `boolean` | True/false values | `true`, `false` |    type = "static"

  include_silent = true

}| `int` | Integer numbers | `42`, `1000` |    properties = jsonencode({



output "importable_attributes" {| `date` | Date values | `"2023-12-25"` |      value = ["technical", "leadership"]

  value = [

    for attr in data.sailpoint_identity_attribute_list.discovery.items :    })

    attr.name if !attr.system

  ]## Searchable vs Non-Searchable  }

}

```}



### 2. Match Existing State- **Searchable** (`searchable = true`): Attributes appear in search interfaces and can be used for filtering identities```

Configure resources to match the existing SailPoint state:

```hcl- **Non-Searchable** (`searchable = false`): Attributes are stored but not indexed for search (useful for sensitive or internal data)

resource "sailpoint_identity_attribute" "imported" {

  name         = "existingAttribute"## Argument Reference

  display_name = "Existing Attribute"  # Match current display name

  type         = "string"               # Match current type## Multi-Value Attributes

  searchable   = true                   # Match current searchability

  The following arguments are supported:

  lifecycle {

    prevent_destroy = true              # Safety during importSet `multi = true` for attributes that can contain multiple values:

    ignore_changes = [standard, system] # Ignore computed fields

  }* `name` - (Required) The technical name of the identity attribute. This cannot be changed after creation.

}

``````hcl* `display_name` - (Optional) The human-readable display name of the identity attribute. Defaults to the value of `name`.



### 3. Gradual Management Transitionresource "sailpoint_identity_attribute" "skills" {* `type` - (Optional) The data type of the identity attribute. Valid values are: `string`, `int`, `boolean`, `date`. Defaults to `string`.

Use lifecycle rules to gradually transition to full management:

```hcl  name         = "skillTags"* `multi` - (Optional) Whether this attribute can have multiple values. Defaults to `false`.

# Phase 1: Import with protection

lifecycle {  display_name = "Skill Tags"* `searchable` - (Optional) Whether this attribute is searchable. Defaults to `false`.

  prevent_destroy = true

  ignore_changes = [standard, system]  type         = "string"* `sources` - (Optional) List of sources that define how this identity attribute gets its values.

}

  multi        = true      # Can contain multiple skill values

# Phase 2: Remove protection, enable management

lifecycle {  searchable   = true### Sources Block

  ignore_changes = [standard, system]

}}



# Phase 3: Full management (remove lifecycle block)```The `sources` block supports:

```



### 4. Validation and Verification

Validate imports using data sources:## Import Existing Attributes* `type` - (Required) The type of source. Valid values include: `rule`, `static`, `connector`.

```hcl

data "sailpoint_identity_attribute" "verify" {* `properties` - (Optional) A JSON string containing the configuration properties for this source.

  name = sailpoint_identity_attribute.imported.name

}Import existing identity attributes into Terraform management:



output "import_verification" {## Attribute Reference

  value = {

    matches = {```bash

      name = sailpoint_identity_attribute.imported.name == data.sailpoint_identity_attribute.verify.name

      type = sailpoint_identity_attribute.imported.type == data.sailpoint_identity_attribute.verify.type# Import by attribute nameIn addition to all arguments above, the following attributes are exported:

    }

  }terraform import sailpoint_identity_attribute.example attributeName

}

```* `standard` - Whether this is a standard identity attribute (computed).



## Import Limitations# Example* `system` - Whether this is a system identity attribute (computed).



### Cannot Importterraform import sailpoint_identity_attribute.cost_center costCenter

- **System attributes**: Built-in SailPoint attributes (e.g., `id`, `created`, `modified`)

- **Standard attributes**: Pre-defined SailPoint attributes unless customized```## Import

- **Silent attributes**: Hidden attributes (use `include_silent=true` to discover)



### Import Considerations

- Attribute names are **case-sensitive**## Common Use CasesIdentity attributes can be imported using their name:

- Only the `name` field is required for import

- Other fields will be computed from SailPoint's current state

- Use `terraform plan` after import to see any configuration drift

### Business Information```bash

## Troubleshooting Imports

- Employee ID, department, cost centerterraform import sailpoint_identity_attribute.cost_center "costCenter"

### Common Issues

- Job title, manager, location```

1. **"Resource not found"**

   - Verify the attribute exists in SailPoint- Organization hierarchy levels

   - Check the exact spelling and case of the attribute name

   - Ensure it's not a system attribute## Notes



2. **"Resource already exists"**### Security & Compliance

   - The resource is already in Terraform state

   - Use `terraform state list` to check existing resources- Security clearance levels- Identity attribute names are immutable after creation

   - Use `terraform state rm` to remove if needed

- Training completion status- System and standard attributes cannot be deleted

3. **Configuration drift after import**

   - Update your `.tf` configuration to match SailPoint's state- Background check dates- Making an attribute searchable requires that `system`, `standard`, and `multi` properties be set to false

   - Use data sources to discover the actual configuration

   - Consider using `ignore_changes` for computed fields- Risk assessment scores- The `properties` field in sources should be valid JSON



### Debug Commands### Technical Attributes

```bash- System access levels

# List all managed resources- Development environment permissions

terraform state list- On-call rotation participation



# Show detailed resource state### Multi-Value Scenarios

terraform state show sailpoint_identity_attribute.example- Professional certifications

- Project assignments

# Remove resource from state (if needed)- Skill tags and competencies

terraform state rm sailpoint_identity_attribute.example- Building access permissions



# Validate configuration## Best Practices

terraform validate

### Naming Conventions

# Check for drift- Use clear, descriptive names

terraform plan- Follow your organization's naming standards

```- Consider using prefixes for categorization (e.g., `hr_`, `security_`)



## Advanced Import Scenarios### Performance Considerations

- Only make attributes `searchable` if they need to be searched

### Bulk Import- Use appropriate data types to optimize storage

Use the provided scripts for bulk import operations:- Consider the impact of multi-value attributes on performance

```bash

# Use the interactive import script### Lifecycle Management

./import.sh- Plan for attribute deprecation and migration

- Use Terraform's lifecycle rules for sensitive changes

# Generate custom bulk import script- Document attribute purposes and data sources

terraform output bulk_import_script > bulk_import.sh

chmod +x bulk_import.sh### Security

./bulk_import.sh- Set `searchable = false` for sensitive information

```- Consider data classification requirements

- Follow principle of least privilege for attribute visibility

### Conditional Import

Import based on discovery results:## Environment-Specific Configurations

```hcl

locals {Use Terraform variables and locals to manage attributes across environments:

  should_import = contains([

    for attr in data.sailpoint_identity_attribute_list.discovery.items : attr.name```hcl

  ], "targetAttribute")variable "environment" {

}  description = "Environment name"

  type        = string

# Only create resource if attribute exists  default     = "dev"

resource "sailpoint_identity_attribute" "conditional" {}

  count = local.should_import ? 1 : 0

  resource "sailpoint_identity_attribute" "env_specific" {

  name = "targetAttribute"  name         = "${var.environment}_testAttribute"

  # ... other configuration  display_name = "Test Attribute (${var.environment})"

}  type         = "string"

```  searchable   = var.environment != "prod"

}

### Import with Validation```

```hcl

resource "sailpoint_identity_attribute" "validated" {## Integration with Other Resources

  name         = "validatedAttribute"

  display_name = "Validated Attribute"Identity attributes can be referenced by other SailPoint resources:

  searchable   = true

  ```hcl

  lifecycle {# Create the attribute

    postcondition {resource "sailpoint_identity_attribute" "department" {

      condition     = self.searchable == true  name         = "department"

      error_message = "Attribute must be searchable after import"  display_name = "Department"

    }  searchable   = true

  }}

}

```# Reference in other configurations

output "department_attribute_name" {

## Related Resources  value = sailpoint_identity_attribute.department.name

}

- `sailpoint_identity_attribute` data source - Query existing attributes```

- `sailpoint_identity_attribute_list` data source - List all attributes

- SailPoint Transform resources - For advanced attribute value transformation## Troubleshooting



## Example Workflows### Common Issues

1. **Name conflicts**: Attribute names must be unique in your SailPoint tenant

For complete workflow examples, see:2. **System attributes**: Cannot manage system-defined attributes through Terraform

- `import-workflow.tf` - Complete import-first workflow3. **Import failures**: Ensure the attribute exists and name matches exactly

- `import-examples.tf` - Comprehensive import scenarios4. **Permission errors**: Verify your SailPoint API credentials have sufficient permissions

- `integration-patterns.tf` - Import in complex environments

### Validation

The import functionality makes it easy to bring existing SailPoint identity attributes under Terraform management, enabling Infrastructure as Code practices for your identity management system.After creating attributes, verify they appear correctly in SailPoint:
- Check the Identity Attributes page in the SailPoint admin console
- Verify searchable attributes appear in search interfaces
- Test multi-value attributes accept multiple values

## Related Resources

- `sailpoint_identity_attribute` data source - Query existing attributes
- `sailpoint_identity_attribute_list` data source - List all attributes
- SailPoint Transform resources - For advanced attribute value transformation

For more examples and advanced patterns, see the other `.tf` files in this directory.