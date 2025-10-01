# Identity Attribute List Data Source Example

This example demonstrates how to use the `sailpoint_identity_attribute_list` data source to retrieve information about all identity attributes in SailPoint Identity Security Cloud (ISC).

## Usage

The data source supports several optional filters:

```hcl
# Get all identity attributes
data "sailpoint_identity_attribute_list" "all" {
  # No filters
}

# Get only searchable identity attributes
data "sailpoint_identity_attribute_list" "searchable_only" {
  searchable_only = true
}

# Include system attributes
data "sailpoint_identity_attribute_list" "with_system" {
  include_system = true
}

# Include silent attributes  
data "sailpoint_identity_attribute_list" "with_silent" {
  include_silent = true
}
```

## Filter Options

- **`include_system`** (bool, optional): Include system-defined attributes (default: false)
- **`include_silent`** (bool, optional): Include silent attributes not shown in UI (default: false)  
- **`searchable_only`** (bool, optional): Only include searchable attributes (default: false)

## Output Structure

The data source returns an `identity_attribute_list` containing all matching identity attributes. Each attribute has the same structure as the single identity attribute data source.

## Use Cases

This data source is useful for:

1. **Discovery**: Finding all available identity attributes in your tenant
2. **Filtering**: Getting subsets of attributes based on specific criteria
3. **Reporting**: Creating summaries and counts of identity attributes
4. **Validation**: Checking if specific attributes exist before referencing them

## Example Filtering

The examples show how to use Terraform's `for` expressions to filter and transform the results:

```hcl
# Get names of only searchable attributes
searchable_names = [
  for attr in data.sailpoint_identity_attribute_list.all.identity_attribute_list : attr.name
  if attr.searchable
]

# Count attributes by type
summary = {
  total = length(data.sailpoint_identity_attribute_list.all.identity_attribute_list)
  searchable = length([for attr in data.sailpoint_identity_attribute_list.all.identity_attribute_list : attr if attr.searchable])
}
```

## Performance Considerations

- Use filters to reduce the amount of data returned when possible
- The `include_system` and `include_silent` filters can significantly increase the result set size
- Consider caching results if you need to reference the same data multiple times