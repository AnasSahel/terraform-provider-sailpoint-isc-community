# Identity Attribute Resource

This resource allows you to create and manage identity attributes in SailPoint Identity Security Cloud (ISC).

## Example Usage

### Minimal Identity Attribute (only name required)

```hcl
resource "sailpoint_identity_attribute" "minimal" {
  name = "costCenter"
  # All other fields have sensible defaults:
  # display_name = "costCenter" (defaults to name)
  # type = "string" (default)
  # multi = false (default)
  # searchable = false (default)
}
```

### Basic Identity Attribute with explicit values

```hcl
resource "sailpoint_identity_attribute" "cost_center" {
  name         = "costCenter"
  display_name = "Cost Center"
  type         = "string"
  multi        = false
  searchable   = true
}
```

### Identity Attribute with Rule Source

```hcl
resource "sailpoint_identity_attribute" "department_code" {
  name         = "departmentCode"
  display_name = "Department Code"
  type         = "string"
  multi        = false
  searchable   = true

  sources {
    type = "rule"
    properties = jsonencode({
      ruleType = "IdentityAttribute"
      ruleName = "Department Code Mapping Rule"
    })
  }
}
```

### Multi-Value Identity Attribute

```hcl
resource "sailpoint_identity_attribute" "skill_tags" {
  name         = "skillTags"
  display_name = "Skill Tags"
  type         = "string"
  multi        = true
  searchable   = true

  sources {
    type = "static"
    properties = jsonencode({
      value = ["technical", "leadership"]
    })
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The technical name of the identity attribute. This cannot be changed after creation.
* `display_name` - (Optional) The human-readable display name of the identity attribute. Defaults to the value of `name`.
* `type` - (Optional) The data type of the identity attribute. Valid values are: `string`, `int`, `boolean`, `date`. Defaults to `string`.
* `multi` - (Optional) Whether this attribute can have multiple values. Defaults to `false`.
* `searchable` - (Optional) Whether this attribute is searchable. Defaults to `false`.
* `sources` - (Optional) List of sources that define how this identity attribute gets its values.

### Sources Block

The `sources` block supports:

* `type` - (Required) The type of source. Valid values include: `rule`, `static`, `connector`.
* `properties` - (Optional) A JSON string containing the configuration properties for this source.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `standard` - Whether this is a standard identity attribute (computed).
* `system` - Whether this is a system identity attribute (computed).

## Import

Identity attributes can be imported using their name:

```bash
terraform import sailpoint_identity_attribute.cost_center "costCenter"
```

## Notes

- Identity attribute names are immutable after creation
- System and standard attributes cannot be deleted
- Making an attribute searchable requires that `system`, `standard`, and `multi` properties be set to false
- The `properties` field in sources should be valid JSON