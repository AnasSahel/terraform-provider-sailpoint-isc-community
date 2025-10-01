# Identity Attribute Data Source Example

This example demonstrates how to use the `sailpoint_identity_attribute` data source to retrieve information about a specific identity attribute in SailPoint Identity Security Cloud (ISC).

## Usage

The data source requires the technical name of the identity attribute you want to retrieve:

```hcl
data "sailpoint_identity_attribute" "example" {
  name = "costCenter"
}
```

## Outputs

The data source provides comprehensive information about the identity attribute:

- **Basic Information**: `name`, `display_name`, `type`
- **Configuration**: `multi`, `searchable`, `system`, `standard`
- **Sources**: Structured list of sources that define how the attribute gets its values

## Sources Structure

The `sources` attribute is a list of objects with the following structure:

```hcl
sources = [
  {
    type       = "rule"           # Type of source (rule, static, accountAttribute, etc.)
    properties = "{...}"          # JSON string containing source-specific configuration
  }
]
```

## Common Identity Attribute Names

Some common identity attribute names you might query:

- `costCenter`
- `department` 
- `division`
- `location`
- `manager`
- `title`
- `employeeNumber`

## Notes

- The `name` parameter is the technical name, not the display name
- System and standard attributes have special behaviors and constraints
- The `properties` field in sources contains JSON configuration specific to each source type